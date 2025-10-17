package httperr

import (
	"errors"
	"fmt"
	"net/http"
)

// Handler is an interface for handling errors and converting them to HTTP responses.
// HandleError should return true if it wrote a response, false otherwise.
type Handler interface {
	HandleError(err error, writer http.ResponseWriter, request *http.Request) (handled bool)
}

// HandlerFunc is an adapter type that allows ordinary functions to be used as error handlers.
// It implements the Handler interface.
type HandlerFunc func(err error, writer http.ResponseWriter, request *http.Request) (handled bool)

// HandleError implements the Handler interface for HandlerFunc.
// It returns false immediately if err is nil, otherwise it calls the underlying function.
func (f HandlerFunc) HandleError(err error, writer http.ResponseWriter, request *http.Request) (handled bool) {
	if err == nil {
		return false
	}
	return f(err, writer, request)
}

// Handle processes an error using the DefaultHandler.
// It returns false if err is nil, otherwise it delegates to DefaultHandler.HandleError.
//
// This is the main entry point for error handling in most cases.
// Use it in handlers that return errors:
//
//	func myHandler(w http.ResponseWriter, r *http.Request) {
//	    err := doSomething()
//	    if httperr.Handle(err, w, r) {
//	        return // Error was handled
//	    }
//	    // Continue with success response
//	}
func Handle(err error, writer http.ResponseWriter, request *http.Request) (handled bool) {
	if err == nil {
		return false
	}
	return DefaultHandler.HandleError(err, writer, request)
}

// HandlePanic processes a panic value recovered by recover() and handles it as an error.
// It converts the panic value to an error using AsError and then calls Handle.
//
// Use this in defer statements to catch panics in HTTP handlers:
//
//	func myHandler(w http.ResponseWriter, r *http.Request) {
//	    defer func() {
//	        if recovered := recover(); recovered != nil {
//	            httperr.HandlePanic(recovered, w, r)
//	        }
//	    }()
//	    // Handler code that might panic
//	}
func HandlePanic(recoverResult any, writer http.ResponseWriter, request *http.Request) (handled bool) {
	return Handle(AsError(recoverResult), writer, request)
}

// ForEachHandler tries multiple error handlers in sequence until one handles the error.
// It returns true if any handler returned true, false otherwise.
//
// This is useful for composing multiple error handling strategies:
//
//	handledAny := httperr.ForEachHandler(err, w, r,
//	    customHandler1,
//	    customHandler2,
//	    httperr.DefaultHandler,
//	)
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
func WriteInternalServerError(err any, writer http.ResponseWriter) {
	message := http.StatusText(http.StatusInternalServerError)
	if DebugShowInternalErrorsInResponse {
		message += fmt.Sprintf(DebugShowInternalErrorsInResponseFormat, err)
	}
	http.Error(writer, message, http.StatusInternalServerError)
}
