package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFile(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	file.Close()

	removeCmd.Run(nil, []string{file.Name()})

	_, err = os.Stat(file.Name())
	assert.True(t, os.IsNotExist(err))
}

func TestRemoveDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	removeCmd.Run(nil, []string{dir})

	_, err = os.Stat(dir)
	assert.True(t, os.IsNotExist(err))
}

func TestRemoveNotExist(t *testing.T) {
	removeCmd.Run(nil, []string{"not-exist"})
}
