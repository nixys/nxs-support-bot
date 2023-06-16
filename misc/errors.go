package misc

import "errors"

var (
	ErrNotFound       = errors.New("entity not found")
	ErrUserCtxExtract = errors.New("can not extract user context in schedule message handler")
	ErrAPIKey         = errors.New("incorrect api key")
	ErrZeroLen        = errors.New("zero len payload")
	ErrForbidden      = errors.New("forbidden")
	ErrMalformedData  = errors.New("malformed data")
	ErrUserNotSet     = errors.New("user not set")
)
