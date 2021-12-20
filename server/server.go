package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
	"io"
	"net"
	"os"
	"qperf-go/common"
	"time"
)

// Run server.
// if proxyAddr is nil, no proxy is used.
func Run(addr net.UDPAddr, createQLog bool, migrateAfter time.Duration, proxyAddr *net.UDPAddr, tlsServerCertFile string, tlsServerKeyFile string, tlsProxyCertFile string, initialCongestionWindow uint32, initialReceiveWindow uint64, maxReceiveWindow uint64) {

	// TODO
	if proxyAddr != nil {
		panic("implement me")
	}

	state := common.State{}

	tracers := make([]logging.Tracer, 0)

	tracers = append(tracers, common.StateTracer{
		State: &state,
	})

	if createQLog {
		tracers = append(tracers, qlog.NewTracer(func(p logging.Perspective, connectionID []byte) io.WriteCloser {
			filename := fmt.Sprintf("server_%x.qlog", connectionID)
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			return common.NewBufferedWriteCloser(bufio.NewWriter(f), f)
		}))
	}

	if initialReceiveWindow > maxReceiveWindow {
		maxReceiveWindow = initialReceiveWindow
	}

	conf := quic.Config{
		Tracer: logging.NewMultiplexedTracer(tracers...),
		IgnoreReceived1RTTPacketsUntilFirstPathMigration: proxyAddr != nil,
		EnableActiveMigration:                            true,
		InitialCongestionWindow:                          initialCongestionWindow,
		InitialStreamReceiveWindow:                       initialReceiveWindow,
		MaxStreamReceiveWindow:                           maxReceiveWindow,
		InitialConnectionReceiveWindow:                   uint64(float64(initialReceiveWindow) * quic.ConnectionFlowControlMultiplier),
		MaxConnectionReceiveWindow:                       uint64(float64(maxReceiveWindow) * quic.ConnectionFlowControlMultiplier),
	}

	//TODO make CLI option
	tlsCert, err := tls.LoadX509KeyPair(tlsServerCertFile, tlsServerKeyFile)
	if err != nil {
		panic(err)
	}

	tlsConf := tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"qperf"},
	}

	listener, err := quic.ListenAddrEarly(addr.String(), &tlsConf, &conf)
	if err != nil {
		panic(err)
	}

	logger := common.DefaultLogger.WithPrefix("") //TODO make cli argument
	logger.SetLogLevel(common.LogLevelInfo)

	logger.Infof("starting server with pid %d, port %d, cc cubic, iw %d", os.Getpid(), addr.Port, conf.InitialCongestionWindow)

	// migrate
	if migrateAfter.Nanoseconds() != 0 {
		go func() {
			time.Sleep(migrateAfter)
			addr, err := listener.MigrateUDPSocket()
			if err != nil {
				panic(err)
			}
			fmt.Printf("migrated to %s\n", addr.String())
		}()
	}

	var nextSessionId uint64 = 0

	for {
		quicSession, err := listener.Accept(context.Background())
		if err != nil {
			panic(err)
		}

		qperfSession := &qperfServerSession{
			session:           quicSession,
			sessionID:         nextSessionId,
			currentRemoteAddr: quicSession.RemoteAddr(),
			logger:            common.DefaultLogger.WithPrefix(fmt.Sprintf("session %d", nextSessionId)),
		}

		go qperfSession.run()
		nextSessionId += 1
	}
}
