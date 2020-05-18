package httperr

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
)

type Handler interface {
	HandleError(err error, writer http.ResponseWriter, request *http.Request) bool
}

type HandlerFunc func(err error, writer http.ResponseWriter, request *http.Request) bool

func (f HandlerFunc) HandleError(err error, writer http.ResponseWriter, request *http.Request) bool {
	return f(err, writer, request)
}

var DefaultHandler Handler = HandlerFunc(DefaultHandlerFunc)

// Handle will call DefaultHandler.HandleError(err, writer, request)
func Handle(err error, writer http.ResponseWriter, request *http.Request) bool {
	if err == nil {
		return false
	}
	return DefaultHandler.HandleError(err, writer, request)
}

// HandlePanic will call DefaultHandler.HandleError(AsError(recoverResult), writer, request)
func HandlePanic(recoverResult interface{}, writer http.ResponseWriter, request *http.Request) bool {
	return Handle(AsError(recoverResult), writer, request)
}

// DefaultHandlerFunc checks err unwraps to a http.Handler and calls its ServeHTTP method if available
// else if err unwraps to os.ErrNotExist or sql.ErrNoRows a 404 Not Found response is written.
// In all other cases a 500 Internal Server Error response is written. with the Error string
// If DebugShowInternalErrorsInResponse is true, then err.Error() message is added to the response.
// If err is nil, then no response is written and the function returns false.
// If an error response was written, then the function returns true.
func DefaultHandlerFunc(err error, writer http.ResponseWriter, request *http.Request) (responseWritten bool) {
	if err == nil {
		return false
	}
	var httperrResponse Response
	switch {
	case errors.As(err, &httperrResponse):
		httperrResponse.ServeHTTP(writer, request)
	case errors.Is(err, os.ErrNotExist):
		http.Error(writer, "Requested file not found", http.StatusNotFound)
	case errors.Is(err, sql.ErrNoRows):
		http.Error(writer, "Requested database row not found", http.StatusNotFound)
	default:
		WriteInternalServerError(err, writer)
	}
	return true
}

func ForEachHandler(err error, writer http.ResponseWriter, request *http.Request, handlers ...Handler) (handledAny bool) {
	if err == nil {
		return false
	}
	for _, handler := range handlers {
		handled := handler.HandleError(err, writer, request)
		handledAny = handledAny || handled
	}
	return handledAny
}
