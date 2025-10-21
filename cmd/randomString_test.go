package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringLength(t *testing.T) {
	src := "abc"
	length := 10
	s := randomString(src, length)
	assert.Equal(t, length, len(s), "Generated string should have the correct length")
}

func TestRandomStringCharset(t *testing.T) {
	src := "abc"
	length := 10
	s := randomString(src, length)
	for _, char := range s {
		assert.True(t, strings.ContainsRune(src, char), "Generated string should only contain characters from the source")
	}
}

func TestRandomStringUniqueness(t *testing.T) {
	src := "abcdefghijklmnopqrstuvwxyz"
	length := 10
	s1 := randomString(src, length)
	s2 := randomString(src, length)
	assert.NotEqual(t, s1, s2, "Consecutive random strings should be different")
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	f()
	w.Close()

	out, _ := io.ReadAll(r)
	return string(out)
}

func TestRandomStringExecute(t *testing.T) {
	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set flags and run the command
	randomStringLength = 12
	randomStringCount = 3
	randomStringCmd.Run(nil, nil)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read the output from the buffer
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Perform assertions
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 3, "Execute should generate the correct number of strings")
	for _, line := range lines {
		assert.Equal(t, 12, len(line), "Each generated string should have the correct length")
	}
}
