package httperr

import (
	"database/sql"
	"net/http"
	"os"
)

var (
	DebugShowInternalErrorsInResponse       bool
	DebugShowInternalErrorsInResponseFormat = "\n%+v"

	DefaultHandler Handler = HandlerFunc(DefaultHandlerImpl)

	// SentinelHandlers is used by DefaultHandlerImpl to map
	// wrapped sentinel errors to corresponding http.Handler.
	// By default os.ErrNotExist and sql.ErrNoRows are mapped
	// to handlers that write a "404 Not Found" response.
	SentinelHandlers = map[error]http.Handler{
		os.ErrNotExist: http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			http.Error(writer, "Requested file not found", http.StatusNotFound)
		}),
		sql.ErrNoRows: http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			http.Error(writer, "Requested database row not found", http.StatusNotFound)
		}),
	}
)

var (
	BadRequest       = New(http.StatusBadRequest)       // 400: RFC 7231, 6.5.1
	Unauthorized     = New(http.StatusUnauthorized)     // 401: RFC 7235, 3.1
	PaymentRequired  = New(http.StatusPaymentRequired)  // 402: RFC 7231, 6.5.2
	Forbidden        = New(http.StatusForbidden)        // 403: RFC 7231, 6.5.3
	NotFound         = New(http.StatusNotFound)         // 404: RFC 7231, 6.5.4
	MethodNotAllowed = New(http.StatusMethodNotAllowed) // 405: RFC 7231, 6.5.5
)

const (
	Handled    = true
	NotHandled = false
)
