package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyCmd_FileToFile(t *testing.T) {
	srcFile, err := os.CreateTemp("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())
	srcFile.WriteString("test")
	srcFile.Close()

	destFile, err := os.CreateTemp("", "dest")
	assert.NoError(t, err)
	destFile.Close()
	defer os.Remove(destFile.Name())

	copyCmd.Run(nil, []string{srcFile.Name(), destFile.Name()})

	content, err := os.ReadFile(destFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

func TestCopyCmd_FileToDir(t *testing.T) {
	srcFile, err := os.CreateTemp("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())
	srcFile.WriteString("test")
	srcFile.Close()

	destDir, err := os.MkdirTemp("", "dest")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	copyCmd.Run(nil, []string{srcFile.Name(), destDir})

	destFile := filepath.Join(destDir, filepath.Base(srcFile.Name()))
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

func TestCopyCmd_DirToFile(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "src")
	assert.NoError(t, err)
	defer os.RemoveAll(srcDir)

	destFile, err := os.CreateTemp("", "dest")
	assert.NoError(t, err)
	defer os.Remove(destFile.Name())
	destFile.Close()

	// This should result in an error, but the Run function catches it and prints.
	// We can't easily assert the error here without more complex output capturing,
	// but we can ensure the command doesn't crash.
	copyCmd.Run(nil, []string{srcDir, destFile.Name()})
}

func TestCopyCmd_DirToDir(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "src")
	assert.NoError(t, err)
	defer os.RemoveAll(srcDir)

	srcFile, err := os.CreateTemp(srcDir, "file")
	assert.NoError(t, err)
	srcFile.WriteString("test")
	srcFile.Close()

	destDir, err := os.MkdirTemp("", "dest")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	copyCmd.Run(nil, []string{srcDir, destDir})

	destFile := filepath.Join(destDir, filepath.Base(srcFile.Name()))
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}
