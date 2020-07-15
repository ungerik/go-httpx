package httperr

import (
	"errors"
	"fmt"
)

// DontLog wraps the passed error
// so that ShouldLog returns true.
//
//   httperr.ShouldLog(httperr.BadRequest) == true
//   httperr.ShouldLog(httperr.DontLog(httperr.BadRequest)) == false
func DontLog(err error) error {
	return errDontLog{err}
}

// ShouldLog checks if the passed error
// has been wrapped with DontLog.
//
//   httperr.ShouldLog(httperr.BadRequest) == true
//   httperr.ShouldLog(httperr.DontLog(httperr.BadRequest)) == false
func ShouldLog(err error) bool {
	var dontLog errDontLog
	return !errors.As(err, &dontLog)
}

type errDontLog struct {
	error
}

func (e errDontLog) Unwrap() error {
	return e.error
}

// AsError converts val to an error by either casting val to error if possible,
// or using its string value or String method as error message,
// or using fmt.Errorf("%+v", val) to format the value as error.
func AsError(val interface{}) error {
	switch x := val.(type) {
	case nil:
		return nil
	case error:
		return x
	case string:
		return errors.New(x)
	case fmt.Stringer:
		return errors.New(x.String())
	default:
		return fmt.Errorf("%+v", val)
	}
}
