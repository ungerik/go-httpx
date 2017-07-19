package respond

import (
	"encoding/xml"
	"net/http"

	"github.com/ungerik/go-httpx/httperr"
)

type XML func(http.ResponseWriter, *http.Request) (response interface{}, err error)

func (handlerFunc XML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

	WriteXML(writer, response)
}

func WriteXML(writer http.ResponseWriter, response interface{}) {
	b, err := EncodeXML(response)
	if err != nil {
		httperr.WriteInternalServerError(err, writer)
		return
	}
	writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	writer.Write([]byte(xml.Header))
	writer.Write(b)
}

func EncodeXML(response interface{}) ([]byte, error) {
	if PrettyPrint {
		return xml.MarshalIndent(response, "", PrettyPrintIndent)
	}
	return xml.Marshal(response)
}
