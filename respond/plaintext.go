package respond

import (
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
	"github.com/ungerik/go-httpx/httperr"
)

type Plaintext func(http.ResponseWriter, *http.Request) (string, error)

func (handlerFunc Plaintext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
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
	writer.Header().Add("Content-Type", contenttype.PlainText)
	writer.Write([]byte(response)) //#nosec G104
}
