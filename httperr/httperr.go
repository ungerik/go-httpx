package httperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// WriteInternalServerError writes err as 500 Internal Server Error reponse.
// If Logger is not nil, then it will be used to log an error message.
// If DebugShowInternalErrorsInResponse is true, then the error message
// will be shown in the response body, else only "Internal Server Error" will be used.
func WriteInternalServerError(err interface{}, writer http.ResponseWriter) {
	message := http.StatusText(http.StatusInternalServerError)
	if DebugShowInternalErrorsInResponse {
		message += fmt.Sprintf(DebugShowInternalErrorsInResponseFormat, err)
	}
	http.Error(writer, message, http.StatusInternalServerError)
}

// WriteAsJSON unmarshals err as JSON and writes it as application/json
// response body using the passed statusCode.
// If err could not be marshalled as JSON, then an internal server error
// will be written instead using WriteInternalServerError with a wrapped erorr message.
func WriteAsJSON(err interface{}, statusCode int, writer http.ResponseWriter) {
	body, e := json.MarshalIndent(err, "", "  ")
	if e != nil {
		e = fmt.Errorf("error while marshalling error value %+v as JSON: %w", err, e)
		WriteInternalServerError(e, writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	writer.Write(body)
}

// AsError converts val to an error by either casting val to error if possible,
// or using its string value or String method as error message,
// or using fmt.Errorf("%+v", val) to format the value as error.
func AsError(val interface{}) error {
	switch x := val.(type) {
	case nil:
		return nil
	case error:
		return x
	case string:
		return errors.New(x)
	case fmt.Stringer:
		return errors.New(x.String())
	default:
		return fmt.Errorf("%+v", val)
	}
}

// DontLog wraps the passed error
// so that ShouldLog returns true.
//
//   httperr.ShouldLog(httperr.BadRequest) == true
//   httperr.ShouldLog(httperr.DontLog(httperr.BadRequest)) == false
func DontLog(err error) error {
	return errDontLog{err}
}

// ShouldLog checks if the passed error
// has been wrapped with DontLog.
//
//   httperr.ShouldLog(httperr.BadRequest) == true
//   httperr.ShouldLog(httperr.DontLog(httperr.BadRequest)) == false
func ShouldLog(err error) bool {
	var dontLog errDontLog
	return !errors.As(err, &dontLog)
}

type errDontLog struct {
	error
}

func (e errDontLog) Unwrap() error {
	return e.error
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
