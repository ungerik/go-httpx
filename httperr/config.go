package httperr

type logger interface {
	Printf(format string, args ...interface{})
}

var (
	Logger logger

	DebugShowInternalErrorsInResponse       bool
	DebugShowInternalErrorsInResponseFormat = "\n%+v"
)
