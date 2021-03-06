package respond

import (
	"net/http"

	"github.com/ungerik/go-httpx/httperr"
)

type Error func(http.ResponseWriter, *http.Request) error

func (handlerFunc Error) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
		}()
	}

	err := handlerFunc(writer, request)

	httperr.Handle(err, writer, request)
}
