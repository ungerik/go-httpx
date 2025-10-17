package respond

import (
	"encoding/json"
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
	"github.com/ungerik/go-httpx/httperr"
)

// JSON is a handler type for functions that return data to be marshaled as JSON.
// The returned data is automatically serialized to JSON and written with
// Content-Type: application/json. Any error is handled by httperr.Handle.
//
// Example:
//
//	http.Handle("/api/users", respond.JSON(func(w http.ResponseWriter, r *http.Request) (any, error) {
//	    users, err := db.GetUsers()
//	    return users, err
//	}))
type JSON func(http.ResponseWriter, *http.Request) (response any, err error)

// ServeHTTP implements http.Handler for JSON.
// It calls the handler function, handles any error, and marshals the response to JSON.
// If CatchPanics is true, panics are recovered and handled as errors.
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

// WriteJSON marshals the response to JSON and writes it with the appropriate content type.
// If marshaling fails, an internal server error is written.
// The response is pretty-printed if PrettyPrint is true.
func WriteJSON(writer http.ResponseWriter, response any) {
	b, err := EncodeJSON(response)
	if err != nil {
		httperr.WriteInternalServerError(err, writer)
		return
	}
	writer.Header().Set("Content-Type", contenttype.JSON)
	writer.Write(b) //#nosec G104
}

// EncodeJSON marshals the response to JSON bytes.
// The response is pretty-printed if PrettyPrint is true.
func EncodeJSON(response any) ([]byte, error) {
	if PrettyPrint {
		return json.MarshalIndent(response, "", PrettyPrintIndent)
	}
	return json.Marshal(response)
}
