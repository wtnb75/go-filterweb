package filterweb

import "errors"

var (
	ErrFilterNotFound      = errors.New("filter not found")
	ErrContentTypeMismatch = errors.New("content type mismatch")
	ErrHTTPRequestFailed   = errors.New("http request failed")
	ErrHTTPStatusNotOK     = errors.New("http status not ok")
	ErrMissingParams       = errors.New("missing parameters")
	ErrReadTemplate        = errors.New("failed to read template")
	ErrEncode              = errors.New("encode error")
	ErrDecode              = errors.New("decode error")
)
