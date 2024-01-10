package respond

import (
	"encoding/json"
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
	"github.com/ungerik/go-httpx/httperr"
)

type JSON func(http.ResponseWriter, *http.Request) (response any, err error)

func (handlerFunc JSON) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
		}()
	}

	response, err := handlerFunc(writer, request)
	if httperr.Handle(err, writer, request) {
		return
	}

	WriteJSON(writer, response)
}

func WriteJSON(writer http.ResponseWriter, response any) {
	b, err := EncodeJSON(response)
	if err != nil {
		httperr.WriteInternalServerError(err, writer)
		return
	}
	writer.Header().Set("Content-Type", contenttype.JSON)
	writer.Write(b) //#nosec G104
}

func EncodeJSON(response any) ([]byte, error) {
	if PrettyPrint {
		return json.MarshalIndent(response, "", PrettyPrintIndent)
	}
	return json.Marshal(response)
}
