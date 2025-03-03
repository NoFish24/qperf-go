package common

import (
	"context"
	"github.com/nofish24/quic-go"
	"github.com/nofish24/quic-go/logging"
)

func NewMultiplexedTracer(tracers ...func(ctx context.Context, perspective logging.Perspective, id quic.ConnectionID) logging.ConnectionTracer) func(ctx context.Context, perspective logging.Perspective, id quic.ConnectionID) logging.ConnectionTracer {
	return func(ctx context.Context, perspective logging.Perspective, id quic.ConnectionID) logging.ConnectionTracer {
		var connectionTracers []logging.ConnectionTracer
		for _, tracer := range tracers {
			connectionTracers = append(connectionTracers, tracer(ctx, perspective, id))
		}
		return logging.NewMultiplexedConnectionTracer(connectionTracers...)
	}
}
