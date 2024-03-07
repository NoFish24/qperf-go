package perf_client

import (
	"crypto/tls"
	"github.com/nofish24/quic-go"
	"github.com/nofish24/quic-go/logging"
	"qperf-go/perf"
)

type Config struct {
	TlsConfig       *tls.Config
	QuicConfig      *quic.Config
	OnStreamSend    func(id quic.StreamID, count logging.ByteCount)
	OnStreamReceive func(id quic.StreamID, count logging.ByteCount)
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
	return c
}
