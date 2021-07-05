// kick:render
package atomic

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/errs"
)

// Atomic atomically writes files by using a temp file.
// When Close is called the temp file is closed and moved to its final destination.
type Atomic struct {
	File    string            // Path to file
	dst     string
	f       *os.File
	written int64
}

// New creates a io.WriteCloser to atomically write files.
func New(a *Atomic) *Atomic {
	if a == nil {
		a = &Atomic{
			dst: dst,
		}
	}
	if a.File == "" {
		panic(`file is not set`)
	}
	a.dst = a.File
	return a
}

// Close closes the temporary file and moves to the destination
func (a *Atomic) Close() error {
	if a.f == nil {
		err := fmt.Errorf("Object is nil")
		if err != nil {
			return err
		}
	}
	a.f.Close()
	err := Move(a.f.Name(), a.dst)
	if err != nil {
		return err
	}
	return nil
}

// Copy Reads until EOF or an error occurs. Data is written to the tempfile
func (a *Atomic) Copy(rdr io.Reader) (written int64, err error) {
	f, err := a.tempfile()
	if err != nil {
		return 0, err
	}
	written, err = io.Copy(f, rdr)
	errs.Panic(err)
	a.written += written
	return written, nil
}

// Write writes bytes to the tempfile
func (a *Atomic) Write(data []byte) (written int, err error) {
	f, err := a.tempfile()
	if err != nil {
		return 0, err
	}
	written, err = f.Write(data)
	if errs.LogF("Can not write to temporary file: %w", err) {
		return written, err
	}
	return written, nil
}

// tempfile returns the *os.File object for the temporary file
func (a *Atomic) tempfile() (*os.File, error) {
	if a.f != nil {
		return a.f, nil
	}
	f, err := ioutil.TempFile("", "")
	if errs.LogF("Can not open temp file: %v", err) {
		return nil, err
	}
	a.f = f
	return a.f, nil
}
