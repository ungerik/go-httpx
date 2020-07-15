package httperr

import (
	"errors"
	"net/http"
)

type Handler interface {
	HandleError(err error, writer http.ResponseWriter, request *http.Request) bool
}

type HandlerFunc func(err error, writer http.ResponseWriter, request *http.Request) bool

func (f HandlerFunc) HandleError(err error, writer http.ResponseWriter, request *http.Request) bool {
	return f(err, writer, request)
}

var DefaultHandler Handler = HandlerFunc(DefaultHandlerImpl)

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

// DefaultHandlerImpl checks if err unwraps to a http.Handler and calls its ServeHTTP method
// else it checks if err wrapped any key in SentinelHandlers and calls ServeHTTP of the http.Handler value.
// In all other cases a 500 Internal Server Error response is written.
// If DebugShowInternalErrorsInResponse is true, then err.Error() message is added to the response.
// If err is nil, then no response is written and the function returns false.
// If an error response was written, then the function returns true.
func DefaultHandlerImpl(err error, writer http.ResponseWriter, request *http.Request) (responseWritten bool) {
	if err == nil {
		return false
	}

	var handler http.Handler
	if errors.As(err, &handler) {
		handler.ServeHTTP(writer, request)
		return true
	}

	for sentinel, handler := range SentinelHandlers {
		if errors.Is(err, sentinel) {
			handler.ServeHTTP(writer, request)
			return true
		}
	}

	WriteInternalServerError(err, writer)
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
