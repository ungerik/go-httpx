package httperr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func WriteInternalServerError(err interface{}, writer http.ResponseWriter) {
	message := http.StatusText(http.StatusInternalServerError)
	if Logger != nil {
		Logger.Printf("%s: %+v", message, err)
	}
	if DebugShowInternalErrorsInResponse {
		message += fmt.Sprintf(DebugShowInternalErrorsInResponseFormat, err)
	}
	http.Error(writer, message, http.StatusInternalServerError)
}

func AsError(p interface{}) error {
	switch x := p.(type) {
	case nil:
		return nil
	case error:
		return x
	case string:
		return errors.New(x)
	case fmt.Stringer:
		return errors.New(x.String())
	}
	return fmt.Errorf("%+v", p)
}

// Response extends the error interface with the http.Handler interface
// to enable errors which can render themselves as HTTP responses.
type Response interface {
	error
	http.Handler
}

func New(statusCode int, statusText ...string) Response {
	return &errWithStatus{
		statusCode: statusCode,
		statusText: strings.Join(statusText, "\n"),
	}
}

type errWithStatus struct {
	statusCode int
	statusText string
}

func (e *errWithStatus) Error() string {
	if e.statusText == "" {
		return http.StatusText(e.statusCode)
	}
	return e.statusText
}

func (e *errWithStatus) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if Logger != nil {
		Logger.Printf("%s", e)
	}
	http.Error(writer, e.Error(), e.statusCode)
}

var (
	BadRequest       = New(http.StatusBadRequest)       // 400: RFC 7231, 6.5.1
	Unauthorized     = New(http.StatusUnauthorized)     // 401: RFC 7235, 3.1
	PaymentRequired  = New(http.StatusPaymentRequired)  // 402: RFC 7231, 6.5.2
	Forbidden        = New(http.StatusForbidden)        // 403: RFC 7231, 6.5.3
	NotFound         = New(http.StatusNotFound)         // 404: RFC 7231, 6.5.4
	MethodNotAllowed = New(http.StatusMethodNotAllowed) // 405: RFC 7231, 6.5.5
)

func Redirect(statusCode int, targetURL string) Response {
	return &redirect{
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

func (r *redirect) Error() string {
	return fmt.Sprintf("%d redirect to %s", r.statusCode, r.targetURL)
}

func (r *redirect) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, r.targetURL, r.statusCode)
}
