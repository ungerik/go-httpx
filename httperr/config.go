package httperr

import "github.com/ungerik/go-httpx"

var (
	Logger httpx.Logger

	DebugShowInternalErrorsInResponse       bool
	DebugShowInternalErrorsInResponseFormat = "\n%+v"
)
