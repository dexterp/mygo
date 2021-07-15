// kick:render
package di

import (
	"fmt"
	"io"
	"log"
	"os"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/errs"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/exit"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/logger"
)

// DI dependency injection struct
type DI struct {
	ExitMode int          // Exit mode one of exit.MNone or exit.Panic
	LogLevel logger.Level // Logger level
	Stderr   io.Writer    // Stderr
	Stdout   io.Writer    // Stdout

	/* Cache objects */
	cacheErrHandler  *errs.Handler
	cacheExitHandler *exit.Handler
	cacheLogFile     *os.File
}

// Options constructor options
type Options struct {
	ExitMode int          // Exit mode one of exit.MNone or exit.Panic
	LogLevel logger.Level // Logger level
	Stderr io.Writer      // Stderr - defaults to os.Stderr
	Stdout io.Writer      // Stdout - defaults to os.Stdout
}

// New DI container constructor
func New(opts Options) *DI {
	di := &DI{
		ExitMode: opts.ExitMode,
		LogLevel: opts.LogLevel,
		Stderr: opts.Stderr,
		Stdout: opts.Stdout,
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
	handler := errs.New(i.MakeExitHandler(), i.MakeLoggerOutput("", ""))
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

// MakeLogFile create a logfile and return the interface
func (i *DI) MakeLogFile(logfile string) *os.File {
	if i.cacheLogFile != nil {
		return i.cacheLogFile
	}
	var (
		f   *os.File
		err error
	)

	fInfo, err := os.Stat(logfile)
	if err != nil && !os.IsNotExist(err) {
		// Simple output because logging is not available
		fmt.Printf(`can not open log file %s: %v`, logfile, err)
	} else if err == nil && fInfo.Size() > 1024*1024*2 {
		// Remove files greater than 2M
		os.Remove(logfile)
	}
	f, err = os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		i.MakeErrorHandler().FatalF(`can not open logfile %s: %v`, logfile, err)
	}
	i.cacheLogFile = f
	return f
}

// MakeLoggerOutput inject logger.OutputIface.
func (i *DI) MakeLoggerOutput(prefix string, logfile string) *logger.Router {
	toStderr := logger.New(i.Stderr, prefix, log.Lmsgprefix, i.LogLevel, i.MakeExitHandler())
	if logfile != "" {
		toFile := logger.New(i.MakeLogFile(logfile), prefix, log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix, i.LogLevel, i.MakeExitHandler())
		return logger.NewRouter(toFile, toStderr)
	}
	return logger.NewRouter(toStderr)
}
