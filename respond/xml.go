package respond

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
				WriteInternalServerError(writer, r)
			}
		}()
	}

	response, err := handlerFunc(writer, request)
	if HandleError(err, writer, request) {
		return
	}

	WriteXML(writer, response)
}

func WriteXML(writer http.ResponseWriter, response interface{}) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := xml.NewEncoder(buf)
	if PrettyPrint {
		encoder.Indent("", PrettyPrintIndent)
	}
	err := encoder.Encode(response)
	if err != nil {
		WriteInternalServerError(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	writer.Write([]byte(xml.Header))
	writer.Write(buf.Bytes())
}
