// kick:render
package errs

import (
	"fmt"
	"io"
	"os"
	"sync"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/exit"
)

// Handler error handling
type Handler struct {
	ex     exit.HandlerIface `validate:"required"` // Exit handler
	mu     *sync.Mutex
	writer io.Writer
}

// New return a *Handler object.
func New(eh *exit.Handler, writer io.Writer) *Handler {
	return &Handler{
		ex:     eh,
		mu:     &sync.Mutex{},
		writer: writer,
	}
}

// Panic will log an error and panic if err is not nil.
func (e *Handler) Panic(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// PanicF will log an error and panic if any argument passed to format is an error
func (e *Handler) PanicF(format string, v ...interface{}) {
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// LogF will log an error if any argument passed to format is an error
func (e *Handler) LogF(format string, v ...interface{}) bool { // nolint
	return e.hasErrPrintf(format, v...)
}

// Fatal will log an error and exit if err is not nil.
func (e *Handler) Fatal(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	e.ex.Exit(255)
}

// FatalF will log an error and exit if any argument passed to fatal is an error
func (e *Handler) FatalF(format string, v ...interface{}) { // nolint
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	e.ex.Exit(255)
}

func (e *Handler) hasErrPrint(err error) bool {
	return err != nil
}

func (e *Handler) hasErrPrintf(format string, v ...interface{}) bool {
	hasError := false
	for _, elm := range v {
		if _, ok := elm.(error); ok {
			hasError = true
			break
		}
	}
	if !hasError {
		return false
	}
	fmt.Fprintf(os.Stderr, format, v...)
	return true
}

// Panic will log an error and panic if err is not nil.
func Panic(err error) {
	e := makeErrors()
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// PanicF will log an error and panic if any argument passed to format is an error
func PanicF(format string, v ...interface{}) {
	e := makeErrors()
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// LogF will log an error if any argument passed to format is an error
func LogF(format string, v ...interface{}) bool { // nolint
	e := makeErrors()
	return e.hasErrPrintf(format, v...)
}

// Fatal will log an error and exit if err is not nil.
func Fatal(err error) {
	e := makeErrors()
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	e.ex.Exit(255)
}

// FatalF will log an error and exit if any argument passed to fatal is an error
func FatalF(format string, v ...interface{}) { // nolint
	e := makeErrors()
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	e.ex.Exit(255)
}

func makeErrors() *Handler {
	eh := &exit.Handler{
		Mode: exit.ExitMode,
	}
	e := &Handler{
		ex: eh,
	}
	return e
}
