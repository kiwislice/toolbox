package copy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF2f(t *testing.T) {
	srcFile, err := os.CreateTemp("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())
	srcFile.WriteString("test")
	srcFile.Close()

	destFile, err := os.CreateTemp("", "dest")
	assert.NoError(t, err)
	destFile.Close()
	defer os.Remove(destFile.Name())

	err = f2f(srcFile.Name(), destFile.Name())
	assert.NoError(t, err)

	content, err := os.ReadFile(destFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

func TestF2d(t *testing.T) {
	srcFile, err := os.CreateTemp("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())
	srcFile.WriteString("test")
	srcFile.Close()

	destDir, err := os.MkdirTemp("", "dest")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	err = f2d(srcFile.Name(), destDir)
	assert.NoError(t, err)

	destFile := filepath.Join(destDir, filepath.Base(srcFile.Name()))
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

func TestD2f(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "src")
	assert.NoError(t, err)
	defer os.RemoveAll(srcDir)

	destFile, err := os.CreateTemp("", "dest")
	assert.NoError(t, err)
	defer os.Remove(destFile.Name())
	destFile.Close()

	err = d2f(srcDir, destFile.Name())
	assert.Error(t, err)
}

func TestD2d(t *testing.T) {
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

	err = d2d(srcDir, destDir)
	assert.NoError(t, err)

	destFile := filepath.Join(destDir, filepath.Base(srcFile.Name()))
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}
