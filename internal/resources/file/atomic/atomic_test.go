// kick:render
package atomic_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/file/atomic"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
)

func TestAtomic(t *testing.T) {
	_ = os.Mkdir(testtools.FixtureDir(), 0755)
	path := filepath.Join(testtools.TempDir(), "atomic.txt")

	_ = os.WriteFile(path, nil, 0644)

	a := io.WriteCloser(atomic.New(&atomic.Atomic{
		File: path,
	}))

	contents := `Hello World`

	w, err := a.Write([]byte(contents))
	if err != nil {
		t.Error(err)
	}
	// Equal length
	assert.Equal(t, w, len(contents))

	// File is empty
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []byte(``), dat)

	// File equals original contents
	a.Close()
	dat, err = ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []byte(contents), dat)
}
