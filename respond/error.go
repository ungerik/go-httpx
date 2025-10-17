package respond

import (
	"net/http"

	"github.com/ungerik/go-httpx/httperr"
)

// Error is a handler type for functions that return only an error.
// The function is responsible for writing the response if there's no error.
// Any returned error is automatically handled by httperr.Handle.
//
// This is useful for handlers that perform actions without returning data:
//
//	http.Handle("/action", respond.Error(func(w http.ResponseWriter, r *http.Request) error {
//	    if err := performAction(); err != nil {
//	        return httperr.Errorf(http.StatusBadRequest, "Action failed: %v", err)
//	    }
//	    w.WriteHeader(http.StatusOK)
//	    return nil
//	}))
type Error func(http.ResponseWriter, *http.Request) error

// ServeHTTP implements http.Handler for Error.
// It calls the handler function and passes any error to httperr.Handle.
// If CatchPanics is true, panics are recovered and handled as errors.
func (handlerFunc Error) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
		}()
	}

	err := handlerFunc(writer, request)

	httperr.Handle(err, writer, request)
}
