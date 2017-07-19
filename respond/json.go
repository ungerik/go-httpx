package respond

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	PrettyPrint       bool
	PrettyPrintIndent = "  "
)

type JSON func(http.ResponseWriter, *http.Request) (response interface{}, err error)

func (handlerFunc JSON) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

	WriteJSON(writer, response)
}

func WriteJSON(writer http.ResponseWriter, response interface{}) {
	b, err := EncodeJSON(response)
	if err != nil {
		WriteInternalServerError(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Write(b)
}

func EncodeJSON(response interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := json.NewEncoder(buf)
	if PrettyPrint {
		encoder.SetIndent("", PrettyPrintIndent)
	}
	err := encoder.Encode(response)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
