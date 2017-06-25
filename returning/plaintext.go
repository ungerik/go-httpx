package returning

import (
	"net/http"
)

type Plaintext func(http.ResponseWriter, *http.Request) (string, error)

func (handlerFunc Plaintext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			if r := recover(); r != nil {
				writeInternalServerError(writer, r)
			}
		}()
	}

	response, err := handlerFunc(writer, request)
	if HandleError(writer, request, err) {
		return
	}

	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.Write([]byte(response))
}

type StaticPlaintext string

func (s StaticPlaintext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.Write([]byte(s))
}
