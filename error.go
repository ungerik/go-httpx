package httpx

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ungerik/IBS-PDC/source/pdc-server/config"
)

var (
	// Debug decides if the text of Error.Internal errors is
	// shown in error responses.
	Debug bool

	// Logger is used to log errors if it is not nil.
	Logger *log.Logger
)

// Error wraps an internal error with a HTTP error code for the response
type Error struct {
	Code     int
	Internal error
}

func (e *Error) Error() string {
	text := strconv.Itoa(e.Code) + ": " + http.StatusText(e.Code)
	if config.Debug && e.Internal != nil {
		text += "\n" + e.Internal.Error()
	}
	return text
}

// NewError creates a new Error on the heap
func NewError(code int, err ...error) *Error {
	var e error
	if len(err) == 1 {
		e = err[0]
	} else if len(err) > 1 {
		panic("NewError err arguments must have length of zero or one")
	}
	return &Error{code, e}
}

// HandleError handles err if it not nil by writing a HTTP response.
// If err is not an instance of Error with an error code,
// then 500 Internal Server Error is used as response code.
func HandleError(err error, writer http.ResponseWriter) {
	if err == nil {
		return
	}
	if Logger != nil {
		Logger.Println(err)
	}
	e, ok := err.(*Error)
	if !ok {
		e = NewError(http.StatusInternalServerError, err)
	}
	http.Error(writer, e.Error(), e.Code)
}
