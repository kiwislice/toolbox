package randomstring

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringLength(t *testing.T) {
	src := "abc"
	length := 10
	s := RandomString(src, length)
	assert.Equal(t, length, len(s), "Generated string should have the correct length")
}

func TestRandomStringCharset(t *testing.T) {
	src := "abc"
	length := 10
	s := RandomString(src, length)
	for _, char := range s {
		assert.True(t, strings.ContainsRune(src, char), "Generated string should only contain characters from the source")
	}
}

func TestRandomStringUniqueness(t *testing.T) {
	src := "abcdefghijklmnopqrstuvwxyz"
	length := 10
	s1 := RandomString(src, length)
	s2 := RandomString(src, length)
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

func TestExecute(t *testing.T) {
	output := captureOutput(func() {
		Execute([]string{"-length", "12", "-count", "3"})
	})

	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 3, "Execute should generate the correct number of strings")
	for _, line := range lines {
		assert.Len(t, line, 12, "Each generated string should have the correct length")
	}
}

func TestPrintDoc(t *testing.T) {
	output := captureOutput(func() {
		PrintDoc()
	})

	assert.Contains(t, output, "產生隨機字串", "PrintDoc should contain the command description")
	assert.Contains(t, output, "-src", "PrintDoc should contain the src flag")
	assert.Contains(t, output, "-length", "PrintDoc should contain the length flag")
	assert.Contains(t, output, "-count", "PrintDoc should contain the count flag")
}
