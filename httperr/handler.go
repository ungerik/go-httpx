package httperr

import (
	"errors"
	"fmt"
	"net/http"
)

type Handler interface {
	HandleError(err error, writer http.ResponseWriter, request *http.Request) (handled bool)
}

type HandlerFunc func(err error, writer http.ResponseWriter, request *http.Request) (handled bool)

func (f HandlerFunc) HandleError(err error, writer http.ResponseWriter, request *http.Request) (handled bool) {
	return f(err, writer, request)
}

// Handle will call DefaultHandler.HandleError(err, writer, request)
func Handle(err error, writer http.ResponseWriter, request *http.Request) (handled bool) {
	if err == nil {
		return false
	}
	return DefaultHandler.HandleError(err, writer, request)
}

// HandlePanic will call DefaultHandler.HandleError(AsError(recoverResult), writer, request)
func HandlePanic(recoverResult interface{}, writer http.ResponseWriter, request *http.Request) (handled bool) {
	return Handle(AsError(recoverResult), writer, request)
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

	if WriteHandler(err, writer, request) {
		return true
	}

	WriteInternalServerError(err, writer)
	return true
}

// WriteHandler checks if err unwraps to a http.Handler and calls its ServeHTTP method
// else it checks if err wrapped any key in SentinelHandlers and calls ServeHTTP of the http.Handler value.
// If an error response was written, then the function returns true.
func WriteHandler(err error, writer http.ResponseWriter, request *http.Request) (responseWritten bool) {
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

	return false
}

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
