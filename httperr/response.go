package httperr

import (
	"fmt"
	"net/http"
	"strings"
)

// Response is an interface that combines error and http.Handler,
// allowing errors to render themselves as HTTP responses.
//
// Errors that implement this interface can be returned from handlers
// and will be automatically rendered with the appropriate status code
// and response body when processed by Handle or DefaultHandler.
type Response interface {
	error
	http.Handler
}

type statusCodeAndText struct {
	statusCode int
	statusText string
}

// New creates a Response error with the given HTTP status code and optional status text.
// If no status text is provided, the standard HTTP status text for the code will be used.
// Multiple status text strings will be joined with newlines.
//
// Example:
//
//	err := httperr.New(http.StatusTeapot) // Uses "I'm a teapot"
//	err := httperr.New(http.StatusBadRequest, "Invalid email format")
func New(statusCode int, statusText ...string) Response {
	return statusCodeAndText{
		statusCode: statusCode,
		statusText: strings.Join(statusText, "\n"),
	}
}

// Errorf creates a Response error with a formatted message using fmt.Sprintf.
// This is similar to New but allows for formatted error messages.
//
// Example:
//
//	err := httperr.Errorf(http.StatusBadRequest, "User %s not found", username)
func Errorf(statusCode int, format string, a ...any) Response {
	return statusCodeAndText{
		statusCode: statusCode,
		statusText: fmt.Sprintf(format, a...),
	}
}

// NewFromResponse creates a Response error from an http.Response.
// It extracts the status code and status text from the response.
// This is useful when proxying or wrapping HTTP responses from other services.
func NewFromResponse(resonse *http.Response) Response {
	return statusCodeAndText{
		statusCode: resonse.StatusCode,
		statusText: resonse.Status,
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
