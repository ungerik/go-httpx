package respond

import (
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
	"github.com/ungerik/go-httpx/httperr"
)

// HTML is a handler type for functions that return HTML content as bytes.
// The returned HTML is automatically written with Content-Type: text/html.
// Any error is handled by httperr.Handle.
//
// Example:
//
//	http.Handle("/page", respond.HTML(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
//	    html := renderTemplate(r)
//	    return html, nil
//	}))
type HTML func(http.ResponseWriter, *http.Request) ([]byte, error)

// ServeHTTP implements http.Handler for HTML.
// It calls the handler function, handles any error, and writes the HTML response.
// If CatchPanics is true, panics are recovered and handled as errors.
func (handlerFunc HTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
		}()
	}

	response, err := handlerFunc(writer, request)
	if httperr.Handle(err, writer, request) {
		return
	}

	WriteHTML(writer, response)
}

// StaticHTML is a handler type for serving static HTML content.
// The HTML string is written with Content-Type: text/html on every request.
//
// Example:
//
//	http.Handle("/static", respond.StaticHTML("<html><body>Hello</body></html>"))
type StaticHTML string

// ServeHTTP implements http.Handler for StaticHTML.
// It writes the static HTML content on every request.
func (s StaticHTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	WriteHTML(writer, []byte(s))
}

// WriteHTML writes the HTML response with the appropriate content type.
func WriteHTML(writer http.ResponseWriter, response []byte) {
	writer.Header().Add("Content-Type", contenttype.HTML)
	writer.Write(response) //#nosec G104
}
