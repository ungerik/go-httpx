package httperr

import (
	"database/sql"
	"net/http"
	"os"
)

var (
	DebugShowInternalErrorsInResponse       bool
	DebugShowInternalErrorsInResponseFormat = "\n%+v"

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
