package httperr

import (
	"fmt"
	"net/http"
	"strings"
)

// Response extends the error interface with the http.Handler interface
// to enable errors which can render themselves as HTTP responses.
type Response interface {
	error
	http.Handler
}

type statusCodeAndText struct {
	statusCode int
	statusText string
}

func New(statusCode int, statusText ...string) Response {
	return statusCodeAndText{
		statusCode: statusCode,
		statusText: strings.Join(statusText, "\n"),
	}
}

func Errorf(statusCode int, format string, a ...interface{}) Response {
	return statusCodeAndText{
		statusCode: statusCode,
		statusText: fmt.Sprintf(format, a...),
	}
}

func (e statusCodeAndText) Error() string {
	if e.statusText == "" {
		return http.StatusText(e.statusCode)
	}
	return e.statusText
}

func (e statusCodeAndText) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, e.Error(), e.statusCode)
}
