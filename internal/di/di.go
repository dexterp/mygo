// kick:render
package di

import (
	"io"
	"os"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/errs"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/exit"
)

// DI dependency injection struct
type DI struct {
	ExitMode int       // Exit mode one of exit.MNone or exit.Panic
	Stderr   io.Writer // Stderr
	Stdout   io.Writer // Stdout

	/* Cache objects */
	cacheErrHandler  *errs.Handler
	cacheExitHandler *exit.Handler
}

// Options constructor options
type Options struct {
	ExitMode int       // Exit mode one of exit.MNone or exit.Panic
	Stderr   io.Writer // Stderr - defaults to os.Stderr
	Stdout   io.Writer // Stdout - defaults to os.Stdout
}

// New DI container constructor
func New(opts Options) *DI {
	di := &DI{
		ExitMode: opts.ExitMode,
		Stderr:   opts.Stderr,
		Stdout:   opts.Stdout,
	}
	if di.Stderr == nil {
		di.Stderr = os.Stderr
	}
	if di.Stdout == nil {
		di.Stdout = os.Stdout
	}
	return di
}

//
// Dependency Injectors
//

// MakeErrorHandler dependency injector
func (i *DI) MakeErrorHandler() *errs.Handler {
	if i.cacheErrHandler != nil {
		return i.cacheErrHandler
	}
	handler := errs.New(i.MakeExitHandler(), i.Stderr)
	i.cacheErrHandler = handler
	return handler
}

// MakeExitHandler dependency injector
func (i *DI) MakeExitHandler() *exit.Handler {
	if i.cacheExitHandler != nil {
		return i.cacheExitHandler
	}
	handler := &exit.Handler{
		Mode: i.ExitMode,
	}
	i.cacheExitHandler = handler
	return handler
}
