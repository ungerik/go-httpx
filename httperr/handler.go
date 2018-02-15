package httperr

import (
	"net/http"

	"github.com/pkg/errors"
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

func DefaultHandlerFunc(err error, writer http.ResponseWriter, request *http.Request) bool {
	if err == nil {
		return false
	}
	if errResponse, ok := errors.Cause(err).(Response); ok {
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
