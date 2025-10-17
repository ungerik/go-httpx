package respond

import (
	"encoding/xml"
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
	"github.com/ungerik/go-httpx/httperr"
)

// XML is a handler type for functions that return data to be marshaled as XML.
// The returned data is automatically serialized to XML and written with
// Content-Type: application/xml. The XML header is automatically prepended.
// Any error is handled by httperr.Handle.
//
// Example:
//
//	http.Handle("/api/user.xml", respond.XML(func(w http.ResponseWriter, r *http.Request) (any, error) {
//	    user, err := db.GetUser(r)
//	    return user, err
//	}))
type XML func(http.ResponseWriter, *http.Request) (response any, err error)

// ServeHTTP implements http.Handler for XML.
// It calls the handler function, handles any error, and marshals the response to XML.
// If CatchPanics is true, panics are recovered and handled as errors.
func (handlerFunc XML) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if CatchPanics {
		defer func() {
			httperr.HandlePanic(recover(), writer, request)
		}()
	}

	response, err := handlerFunc(writer, request)
	if httperr.Handle(err, writer, request) {
		return
	}

	WriteXML(writer, response)
}

// WriteXML marshals the response to XML and writes it with the appropriate content type.
// The XML header is automatically prepended. If marshaling fails, an internal server error is written.
// The response is pretty-printed if PrettyPrint is true.
func WriteXML(writer http.ResponseWriter, response any) {
	b, err := EncodeXML(response)
	if err != nil {
		httperr.WriteInternalServerError(err, writer)
		return
	}
	writer.Header().Set("Content-Type", contenttype.XML)
	writer.Write([]byte(xml.Header)) //#nosec G104
	writer.Write(b)                  //#nosec G104
}

// EncodeXML marshals the response to XML bytes.
// The response is pretty-printed if PrettyPrint is true.
func EncodeXML(response any) ([]byte, error) {
	if PrettyPrint {
		return xml.MarshalIndent(response, "", PrettyPrintIndent)
	}
	return xml.Marshal(response)
}
