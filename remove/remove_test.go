package remove

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFile(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	file.Close()

	cmd, err := newRemoveCmd([]string{file.Name()})
	assert.NoError(t, err)

	cmd.Run()

	_, err = os.Stat(file.Name())
	assert.True(t, os.IsNotExist(err))
}

func TestRemoveDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	cmd, err := newRemoveCmd([]string{dir})
	assert.NoError(t, err)

	cmd.Run()

	_, err = os.Stat(dir)
	assert.True(t, os.IsNotExist(err))
}

func TestRemoveNotExist(t *testing.T) {
	cmd, err := newRemoveCmd([]string{"not-exist"})
	assert.NoError(t, err)

	cmd.Run()
}
