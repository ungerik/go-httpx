package httpx

type Logger interface {
	Printf(format string, args ...any)
}
