package respond

import (
	"net/http"
)

type HTML func(http.ResponseWriter, *http.Request) ([]byte, error)

func (handlerFunc HTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			if r := recover(); r != nil {
				WriteInternalServerError(writer, r)
			}
		}()
	}

	response, err := handlerFunc(writer, request)
	if HandleError(err, writer, request) {
		return
	}

	WriteHTML(writer, response)
}

type StaticHTML string

func (s StaticHTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	WriteHTML(writer, []byte(s))
}

func WriteHTML(writer http.ResponseWriter, response []byte) {
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.Write(response)
}
