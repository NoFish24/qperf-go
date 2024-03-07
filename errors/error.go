package errors

import "github.com/nofish24/quic-go"

const (
	NoError           = quic.ApplicationErrorCode(0)
	InternalErrorCode = quic.ApplicationErrorCode(1)
)
