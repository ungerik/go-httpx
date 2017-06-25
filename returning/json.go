package returning

import (
	"encoding/json"
	"net/http"
)

var (
	PrettyPrintResponses bool
	PrettyPrintIndent    = "  "
)

type JSON func(http.ResponseWriter, *http.Request) (response interface{}, err error)

func (handlerFunc JSON) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(writer)
	if PrettyPrintResponses {
		encoder.SetIndent("", PrettyPrintIndent)
	}
	err = encoder.Encode(response)
	if err != nil {
		writeInternalServerError(writer, err)
	}
}
