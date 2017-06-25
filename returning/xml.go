package returning

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

type XML func(http.ResponseWriter, *http.Request) (response interface{}, err error)

func (handlerFunc XML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := xml.NewEncoder(buf)
	if PrettyPrintResponses {
		encoder.Indent("", PrettyPrintIndent)
	}
	err = encoder.Encode(response)
	if err != nil {
		writeInternalServerError(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	writer.Write([]byte(xml.Header))
	writer.Write(buf.Bytes())
}
