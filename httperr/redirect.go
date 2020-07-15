package httperr

import (
	"fmt"
	"net/http"
)

func Redirect(statusCode int, targetURL string) Response {
	return redirect{
		statusCode: statusCode,
		targetURL:  targetURL,
	}
}

func TemporaryRedirect(targetURL string) Response {
	return Redirect(http.StatusTemporaryRedirect, targetURL)
}

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
