package client

import (
	"crypto/tls"
	"github.com/nofish24/quic-go"
	"github.com/nofish24/quic-go/logging"
	qlog2 "qperf-go/common/qlog"
	"qperf-go/perf"
	"runtime/debug"
	"time"
)

const (
	DefaultProbeTime      = 10 * time.Second
	MaxProbeTime          = 10 * time.Hour
	DefaultReportInterval = 1 * time.Second
	DefaultQlogTitle      = "qperf"
)

func getDefaultQlogCodeVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Version
}

type Config struct {
	TimeToFirstByteOnly   bool
	ProbeTime             time.Duration
	ReportInterval        time.Duration
	Use0RTT               bool
	LogPrefix             string
	SendInfiniteStream    bool
	ReceiveInfiniteStream bool
	SendDatagram          bool
	ReceiveDatagram       bool
	// output path of qlog file. {odcid} is substituted.
	QlogPathTemplate          string
	QlogConfig                *qlog2.Config
	RemoteAddress             string
	TlsConfig                 *tls.Config
	ReportLostPackets         bool
	ReportMaxRTT              bool
	QuicConfig                *quic.Config
	ReconnectOnTimeoutOrReset bool
	RequestLength             uint64
	ResponseLength            uint64
	RequestInterval           time.Duration
	Deadline                  time.Duration
	ResponseDelay             time.Duration
	//ROSA
	EdgeAddr   string
	ClientAddr string
}

func (c *Config) Populate() *Config {
	if c == nil {
		c = &Config{}
	}
	if c.TlsConfig == nil {
		c.TlsConfig = &tls.Config{}
	}
	if c.TlsConfig.NextProtos == nil {
		c.TlsConfig.NextProtos = []string{perf.ALPN}
	}
	if c.QlogConfig == nil {
		c.QlogConfig = &qlog2.Config{}
		c.QlogConfig.VantagePoint = logging.PerspectiveClient
	}
	if c.QlogConfig.Title == "" {
		c.QlogConfig.Title = DefaultQlogTitle
	}
	if c.QlogConfig.CodeVersion == "" {
		c.QlogConfig.CodeVersion = getDefaultQlogCodeVersion()
	}
	c.QlogConfig.Populate()
	if c.QuicConfig == nil {
		c.QuicConfig = &quic.Config{}
		c.QuicConfig.EnableDatagrams = true
	}
	return c
}
