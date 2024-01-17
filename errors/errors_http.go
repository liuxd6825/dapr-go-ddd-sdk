package errors

import "net/http"

var (
	ErrNotFound            = New("not found")
	ErrHttpBodyNotAllowed  = http.ErrBodyNotAllowed
	ErrHttpHijacked        = http.ErrHijacked
	ErrHttpContentLength   = http.ErrContentLength
	ErrHttpWriteAfterFlush = http.ErrWriteAfterFlush
)
