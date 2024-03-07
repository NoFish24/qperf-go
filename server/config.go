package server

import (
	"crypto/tls"
	"github.com/nofish24/quic-go"
	"github.com/nofish24/quic-go/logging"
	qlog2 "qperf-go/common/qlog"
	"qperf-go/perf"
	"runtime/debug"
)

const (
	DefaultQlogTitle = "qperf"
)

func getDefaultQlogCodeVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Version
}

type Config struct {
	// output path of qlog file. {odcid} is substituted.
	QlogPathTemplate string
	QlogConfig       *qlog2.Config
	TlsConfig        *tls.Config
	QuicConfig       *quic.Config
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
		c.QlogConfig.VantagePoint = logging.PerspectiveServer
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
		c.QuicConfig.Allow0RTT = true
		c.QuicConfig.EnableDatagrams = true
	}
	return c
}
