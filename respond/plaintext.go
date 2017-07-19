package respond

import (
	"net/http"

	"github.com/ungerik/go-httpx/httperr"
)

type Plaintext func(http.ResponseWriter, *http.Request) (string, error)

func (handlerFunc Plaintext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			if r := recover(); r != nil {
				httperr.WriteInternalServerError(r, writer)
			}
		}()
	}

	response, err := handlerFunc(writer, request)
	if httperr.Handle(err, writer, request) {
		return
	}

	WritePlaintext(writer, response)
}

type StaticPlaintext string

func (s StaticPlaintext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	WritePlaintext(writer, string(s))
}

func WritePlaintext(writer http.ResponseWriter, response string) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.Write([]byte(response))
}
