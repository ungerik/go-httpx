package returning

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	CatchPanics                       bool
	DebugShowInternalErrorsInResponse bool
)

type Error func(http.ResponseWriter, *http.Request) error

func (handlerFunc Error) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			if r := recover(); r != nil {
				writeInternalServerError(writer, r)
			}
		}()
	}
	err := handlerFunc(writer, request)
	HandleError(writer, request, err)
}

func writeInternalServerError(writer http.ResponseWriter, err interface{}) {
	message := http.StatusText(http.StatusInternalServerError)
	if DebugShowInternalErrorsInResponse {
		message += fmt.Sprintf(": %+v", err)
	}
	http.Error(writer, message, http.StatusInternalServerError)
}

func HandleError(writer http.ResponseWriter, request *http.Request, err error) bool {
	switch e := err.(type) {
	case nil:
		return false
	case ErrorHandler:
		e.ServeHTTP(writer, request)
		return true
	default:
		writeInternalServerError(writer, err)
		return true
	}
}

type ErrorHandler interface {
	error
	http.Handler
}

func NewError(statusCode int, statusText ...string) ErrorHandler {
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
	http.Error(writer, e.Error(), e.statusCode)
}

var (
	BadRequest       = NewError(http.StatusBadRequest)       // 400: RFC 7231, 6.5.1
	Unauthorized     = NewError(http.StatusUnauthorized)     // 401: RFC 7235, 3.1
	PaymentRequired  = NewError(http.StatusPaymentRequired)  // 402: RFC 7231, 6.5.2
	Forbidden        = NewError(http.StatusForbidden)        // 403: RFC 7231, 6.5.3
	NotFound         = NewError(http.StatusNotFound)         // 404: RFC 7231, 6.5.4
	MethodNotAllowed = NewError(http.StatusMethodNotAllowed) // 405: RFC 7231, 6.5.5
)

func NewRedirect(statusCode int, targetURL string) ErrorHandler {
	return &redirect{
		statusCode: statusCode,
		targetURL:  targetURL,
	}
}

func NewTemporaryRedirect(targetURL string) ErrorHandler {
	return NewRedirect(http.StatusTemporaryRedirect, targetURL)
}

func NewPermanentRedirect(targetURL string) ErrorHandler {
	return NewRedirect(http.StatusPermanentRedirect, targetURL)
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
