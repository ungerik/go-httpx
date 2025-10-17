// Package httpx provides useful extensions around Go's net/http package
// including error handling, response writers, and graceful server shutdown.
//
// The package consists of several sub-packages:
//   - httperr: HTTP error handling and error-to-response conversion
//   - respond: Simplified response writing for JSON, XML, HTML, and plain text
//   - contenttype: Constants for common MIME content types
//   - calling: Function calling utilities with string arguments
package httpx

// Logger is an interface for logging messages.
// It is used by GracefulShutdownServerOnSignal to log signals and errors.
// The standard library's log.Logger implements this interface.
type Logger interface {
	Printf(format string, args ...any)
}
