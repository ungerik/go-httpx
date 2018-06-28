package httperr

import (
	"net/http"
)

type Handler interface {
	HandleError(err error, writer http.ResponseWriter, request *http.Request) bool
}

type HandlerFunc func(err error, writer http.ResponseWriter, request *http.Request) bool

func (f HandlerFunc) HandleError(err error, writer http.ResponseWriter, request *http.Request) bool {
	return f(err, writer, request)
}

var DefaultHandler Handler = HandlerFunc(DefaultHandlerFunc)

func Handle(err error, writer http.ResponseWriter, request *http.Request) bool {
	return DefaultHandler.HandleError(err, writer, request)
}

func RecoverAndHandlePanic(writer http.ResponseWriter, request *http.Request) bool {
	return DefaultHandler.HandleError(AsError(recover()), writer, request)
}

func DefaultHandlerFunc(err error, writer http.ResponseWriter, request *http.Request) bool {
	if err == nil {
		return false
	}
	if errResponse, ok := AsResponse(err); ok {
		errResponse.ServeHTTP(writer, request)
	} else {
		WriteInternalServerError(err, writer)
	}
	return true
}

func ForEachHandler(err error, writer http.ResponseWriter, request *http.Request, handlers ...Handler) (handledAny bool) {
	for _, handler := range handlers {
		handled := handler.HandleError(err, writer, request)
		handledAny = handledAny || handled
	}
	return handledAny
}
