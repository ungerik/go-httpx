package returning

import (
	"net/http"
)

type HTML func(http.ResponseWriter, *http.Request) ([]byte, error)

func (handlerFunc HTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			if r := recover(); r != nil {
				writeInternalServerError(writer, r)
			}
		}()
	}

	response, err := handlerFunc(writer, request)
	if err != nil {
		if errHandler, ok := err.(ErrorHandler); ok {
			errHandler.ServeHTTP(writer, request)
		} else {
			writeInternalServerError(writer, err)
		}
		return
	}

	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.Write(response)
}

type StaticHTML string

func (s StaticHTML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.Write([]byte(s))
}
