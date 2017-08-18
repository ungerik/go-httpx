package httperr

import "log"

var (
	DebugShowInternalErrorsInResponse       bool
	DebugShowInternalErrorsInResponseFormat = "\n%+v"

	Logger *log.Logger
)
