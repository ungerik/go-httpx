package respond

import (
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
	"github.com/ungerik/go-httpx/httperr"
)

// Plaintext is a handler type for functions that return plain text content.
// The returned string is automatically written with Content-Type: text/plain.
// Any error is handled by httperr.Handle.
//
// Example:
//
//	http.Handle("/status", respond.Plaintext(func(w http.ResponseWriter, r *http.Request) (string, error) {
//	    status := checkSystemStatus()
//	    return status, nil
//	}))
type Plaintext func(http.ResponseWriter, *http.Request) (string, error)

// ServeHTTP implements http.Handler for Plaintext.
// It calls the handler function, handles any error, and writes the plain text response.
// If CatchPanics is true, panics are recovered and handled as errors.
func (handlerFunc Plaintext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
		}()
	}

	response, err := handlerFunc(writer, request)
	if httperr.Handle(err, writer, request) {
		return
	}

	WritePlaintext(writer, response)
}

// StaticPlaintext is a handler type for serving static plain text content.
// The text string is written with Content-Type: text/plain on every request.
//
// Example:
//
//	http.Handle("/version", respond.StaticPlaintext("v1.2.3"))
type StaticPlaintext string

// ServeHTTP implements http.Handler for StaticPlaintext.
// It writes the static plain text content on every request.
func (s StaticPlaintext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	WritePlaintext(writer, string(s))
}

// WritePlaintext writes the plain text response with the appropriate content type.
func WritePlaintext(writer http.ResponseWriter, response string) {
	writer.Header().Add("Content-Type", contenttype.PlainText)
	writer.Write([]byte(response)) //#nosec G104
}
