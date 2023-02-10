package respond

import (
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
	"github.com/ungerik/go-httpx/httperr"
)

type HTML func(http.ResponseWriter, *http.Request) ([]byte, error)

func (handlerFunc HTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
		}()
	}

	response, err := handlerFunc(writer, request)
	if httperr.Handle(err, writer, request) {
		return
	}

	WriteHTML(writer, response)
}

type StaticHTML string

func (s StaticHTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	WriteHTML(writer, []byte(s))
}

func WriteHTML(writer http.ResponseWriter, response []byte) {
	writer.Header().Add("Content-Type", contenttype.HTML)
	writer.Write(response) //#nosec G104
}
