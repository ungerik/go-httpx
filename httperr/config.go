// Package httperr provides HTTP error handling functionality that allows
// errors to be used as HTTP responses with status codes and custom messages.
//
// The package enables a clean error handling pattern where errors can be
// returned from handlers and automatically converted to appropriate HTTP responses.
//
// Key features:
//   - Pre-defined error responses for common HTTP status codes
//   - Custom error responses with JSON support
//   - HTTP redirects as error values
//   - Sentinel error mapping (e.g., os.ErrNotExist -> 404)
//   - Error logging control with DontLog wrapper
//   - Panic recovery and conversion to errors
//
// Example usage:
//
//	func handler(w http.ResponseWriter, r *http.Request) error {
//	    if !isValid(r) {
//	        return httperr.BadRequest
//	    }
//	    user, err := getUser(r)
//	    if err != nil {
//	        return err // Automatically handled by httperr.Handle()
//	    }
//	    return nil
//	}
package httperr

import (
	"database/sql"
	"net/http"
	"os"
)

var (
	// DebugShowInternalErrorsInResponse controls whether internal server errors
	// (500) include the actual error message in the response body.
	// This should be set to true in development and false in production.
	DebugShowInternalErrorsInResponse bool

	// DebugShowInternalErrorsInResponseFormat is the format string used to
	// append error details to internal server error responses when
	// DebugShowInternalErrorsInResponse is true.
	DebugShowInternalErrorsInResponseFormat = "\n%+v"

	// DefaultHandler is the error handler used by Handle() and HandlePanic().
	// It can be replaced with a custom handler to change the default error
	// handling behavior globally.
	DefaultHandler Handler = HandlerFunc(DefaultHandlerImpl)

	// SentinelHandlers maps sentinel errors to corresponding http.Handler implementations.
	// When an error that wraps any key in this map is handled by DefaultHandlerImpl,
	// the corresponding handler's ServeHTTP method will be called.
	//
	// By default, the following mappings are configured:
	//   - os.ErrNotExist -> 404 Not Found (file not found)
	//   - sql.ErrNoRows  -> 404 Not Found (database row not found)
	//
	// You can add custom mappings for your own sentinel errors:
	//
	//	httperr.SentinelHandlers[myapp.ErrRateLimited] = http.HandlerFunc(
	//	    func(w http.ResponseWriter, _ *http.Request) {
	//	        http.Error(w, "Rate limited", http.StatusTooManyRequests)
	//	    },
	//	)
	SentinelHandlers = map[error]http.Handler{
		os.ErrNotExist: http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			http.Error(writer, "Requested file not found", http.StatusNotFound)
		}),
		sql.ErrNoRows: http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			http.Error(writer, "Requested database row not found", http.StatusNotFound)
		}),
	}
)

// Pre-defined HTTP error responses for common status codes.
// These can be returned directly from handlers or wrapped with additional context.
var (
	BadRequest       = New(http.StatusBadRequest)       // 400: RFC 7231, 6.5.1
	Unauthorized     = New(http.StatusUnauthorized)     // 401: RFC 7235, 3.1
	PaymentRequired  = New(http.StatusPaymentRequired)  // 402: RFC 7231, 6.5.2
	Forbidden        = New(http.StatusForbidden)        // 403: RFC 7231, 6.5.3
	NotFound         = New(http.StatusNotFound)         // 404: RFC 7231, 6.5.4
	MethodNotAllowed = New(http.StatusMethodNotAllowed) // 405: RFC 7231, 6.5.5
)

// Constants for clarity when returning boolean values from error handlers.
const (
	Handled    = true  // Indicates that an error was handled and a response was written
	NotHandled = false // Indicates that an error was not handled
)
