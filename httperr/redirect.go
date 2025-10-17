package httperr

import (
	"fmt"
	"net/http"
)

// Redirect creates a Response error that performs an HTTP redirect to targetURL
// with the specified status code (typically 301, 302, 307, or 308).
//
// Example:
//
//	return httperr.Redirect(http.StatusMovedPermanently, "/new-location")
func Redirect(statusCode int, targetURL string) Response {
	return redirect{
		statusCode: statusCode,
		targetURL:  targetURL,
	}
}

// TemporaryRedirect creates a Response error that performs a 307 Temporary Redirect
// to the specified targetURL. This preserves the request method and body.
//
// Example:
//
//	return httperr.TemporaryRedirect("/login")
func TemporaryRedirect(targetURL string) Response {
	return Redirect(http.StatusTemporaryRedirect, targetURL)
}

// PermanentRedirect creates a Response error that performs a 308 Permanent Redirect
// to the specified targetURL. This preserves the request method and body.
//
// Example:
//
//	return httperr.PermanentRedirect("/new-permanent-location")
func PermanentRedirect(targetURL string) Response {
	return Redirect(http.StatusPermanentRedirect, targetURL)
}

type redirect struct {
	statusCode int
	targetURL  string
}

func (r redirect) Error() string {
	return fmt.Sprintf("%d redirect to %s", r.statusCode, r.targetURL)
}

func (r redirect) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, r.targetURL, r.statusCode)
}
