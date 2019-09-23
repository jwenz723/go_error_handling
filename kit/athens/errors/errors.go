package errors

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"runtime"
)

// Kind enums
const (
	KindNotFound       = http.StatusNotFound
	KindBadRequest     = http.StatusBadRequest
	KindUnexpected     = http.StatusInternalServerError
	KindAlreadyExists  = http.StatusConflict
	KindRateLimit      = http.StatusTooManyRequests
	KindNotImplemented = http.StatusNotImplemented
	KindRedirect       = http.StatusMovedPermanently
)

// Error is an Athens system error.
// It carries information and behavior
// as to what caused this error so that
// callers can implement logic around it.
type Error struct {
	// Kind categories Athens errors into a smaller
	// subset of errors. This way we can generalize
	// what an error really is: such as "not found",
	// "bad request", etc. The official categories
	// are HTTP status code but the ones we use are
	// imported into this package.
	Kind       int
	Op         Op
	CustomerID C
	OrderID    O
	Err        error
	Severity   level.Value
}

// Error returns the underlying error's
// string message. The logger takes care
// of filling out the stack levels and
// extra information.
func (e Error) Error() string {
	return e.Err.Error()
}

// Is is a shorthand for checking an error against a kind.
func Is(err error, kind int) bool {
	if err == nil {
		return false
	}
	return Kind(err) == kind
}

// Op describes any independent function or
// method in Athens. A series of operations
// forms a more readable stack trace.
type Op string

func (o Op) String() string {
	return string(o)
}

// C represents a customerID
type C string

// O represents an orderID
type O string

// E is a helper function to construct an Error type
// Operation always comes first, module path and version
// come second, they are optional. Args must have at least
// an error or a string to describe what exactly went wrong.
// You can optionally pass a Logrus severity to indicate
// the log level of an error based on the context it was constructed in.
func E(op Op, args ...interface{}) error {
	e := Error{Op: op}
	if len(args) == 0 {
		msg := "errors.E called with 0 args"
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = fmt.Sprintf("%v - %v:%v", msg, file, line)
		}
		e.Err = errors.New(msg)
	}
	for _, a := range args {
		switch a := a.(type) {
		case error:
			e.Err = a
		case string:
			e.Err = errors.New(a)
		case C:
			e.CustomerID = a
		case O:
			e.OrderID = a
		case level.Value:
			e.Severity = a
		case int:
			e.Kind = a
		}
	}
	if e.Err == nil {
		e.Err = errors.New(KindText(e))
	}
	return e
}

// Severity returns the log level of an error
// if none exists, then the level is Error because
// it is an unexpected.
func Severity(err error) level.Value {
	e, ok := err.(Error)
	if !ok {
		return level.ErrorValue()
	}

	// if there's no severity (0 is Panic level in logrus
	// which we should not use since cloud providers only have
	// debug, info, warn, and error) then look for the
	// child's severity.
	// TODO: figure out how to do this check using go-kit level
	//if e.Severity < level.ErrorValue() {
	//	return Severity(e.Err)
	//}

	return e.Severity
}

// Expect is a helper that returns an Info level
// if the error has the expected kind, otherwise
// it returns an Error level.
func Expect(err error, kinds ...int) level.Value {
	for _, kind := range kinds {
		if Kind(err) == kind {
			return level.InfoValue()
		}
	}
	return level.ErrorValue()
}

func CustomerID(err error) C {
	e, ok := err.(Error)
	if !ok {
		return ""
	}

	if e.CustomerID != "" {
		return e.CustomerID
	}

	return CustomerID(e.Err)
}

func OrderID(err error) O {
	e, ok := err.(Error)
	if !ok {
		return ""
	}

	if e.OrderID != "" {
		return e.OrderID
	}

	return OrderID(e.Err)
}

// Kind recursively searches for the
// first error kind it finds.
func Kind(err error) int {
	e, ok := err.(Error)
	if !ok {
		return KindUnexpected
	}

	if e.Kind != 0 {
		return e.Kind
	}

	return Kind(e.Err)
}

// KindText returns a friendly string
// of the Kind type. Since we use http
// status codes to represent error kinds,
// this method just deferrs to the net/http
// text representations of statuses.
func KindText(err error) string {
	return http.StatusText(Kind(err))
}

// Ops aggregates the error's operation
// with all the embedded errors' operations.
// This way you can construct a queryable
// stack trace.
func Ops(err Error) []Op {
	ops := []Op{err.Op}
	for {
		embeddedErr, ok := err.Err.(Error)
		if !ok {
			break
		}

		ops = append(ops, embeddedErr.Op)
		err = embeddedErr
	}

	return ops
}

func OpsText(err Error) string {
	ops := string(err.Op)
	for {
		embeddedErr, ok := err.Err.(Error)
		if !ok {
			break
		}

		ops += fmt.Sprintf(": %s", embeddedErr.Op)
		err = embeddedErr
	}

	return ops
}