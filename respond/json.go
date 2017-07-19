package respond

import (
	"encoding/json"
	"net/http"

	"github.com/ungerik/go-httpx/httperr"
)

type JSON func(http.ResponseWriter, *http.Request) (response interface{}, err error)

func (handlerFunc JSON) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer httperr.Handle(httperr.Recover(), writer, request)
	}

	response, err := handlerFunc(writer, request)
	if httperr.Handle(err, writer, request) {
		return
	}

	WriteJSON(writer, response)
}

func WriteJSON(writer http.ResponseWriter, response interface{}) {
	b, err := EncodeJSON(response)
	if err != nil {
		httperr.WriteInternalServerError(err, writer)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Write(b)
}

func EncodeJSON(response interface{}) ([]byte, error) {
	if PrettyPrint {
		return json.MarshalIndent(response, "", PrettyPrintIndent)
	}
	return json.Marshal(response)
}
