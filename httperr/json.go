package httperr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ungerik/go-httpx/contenttype"
)

// WriteAsJSON unmarshals err as JSON and writes it as application/json
// response body using the passed statusCode.
// If err could not be marshalled as JSON, then an internal server error
// will be written instead using WriteInternalServerError with a wrapped erorr message.
func WriteAsJSON(err any, statusCode int, writer http.ResponseWriter) {
	body, e := json.MarshalIndent(err, "", "  ")
	if e != nil {
		e = fmt.Errorf("can't marshall error of type %T as JSON because: %w", err, e)
		WriteInternalServerError(e, writer)
		return
	}

	writer.Header().Set("Content-Type", contenttype.JSON)
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	writer.Write(body) //#nosec G104
}

// JSON returns a Response error that will respond with
// the passed statusCode, the content type application/json
// and the passed body marshalled as JSON.
// Pass a json.RawMessage as body if the error
// message is already in JSON format.
func JSON(statusCode int, body any) Response {
	return statusCodeAndJSON{
		statusCode: statusCode,
		body:       body,
	}
}

type statusCodeAndJSON struct {
	statusCode int
	body       any
}

func (e statusCodeAndJSON) Error() string {
	body, err := json.MarshalIndent(e.body, "", "  ")
	if err != nil {
		body = []byte(fmt.Sprintf("can't marshall error of type %T as JSON because: %s", e.body, err))
	}
	return fmt.Sprintf("%d: %s", e.statusCode, body)
}

func (e statusCodeAndJSON) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	WriteAsJSON(e.body, e.statusCode, writer)
}
