// 該檔案為 randomString command 的單元測試
package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRandomStringLength 測試產生的隨機字串長度是否正確
func TestRandomStringLength(t *testing.T) {
	src := "abc"
	length := 10
	s := randomString(src, length)
	assert.Equal(t, length, len(s), "產生的字串應具有正確的長度")
}

// TestRandomStringCharset 測試產生的隨機字串是否只包含來源字元
func TestRandomStringCharset(t *testing.T) {
	src := "abc"
	length := 10
	s := randomString(src, length)
	for _, char := range s {
		assert.True(t, strings.ContainsRune(src, char), "產生的字串應只包含來源字元")
	}
}

// TestRandomStringUniqueness 測試連續產生的隨機字串是否不同
func TestRandomStringUniqueness(t *testing.T) {
	src := "abcdefghijklmnopqrstuvwxyz"
	length := 10
	s1 := randomString(src, length)
	s2 := randomString(src, length)
	assert.NotEqual(t, s1, s2, "連續的隨機字串應不同")
}

// captureOutput 是一個輔助函式，用於捕獲函式的標準輸出
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

// TestRandomStringExecute 測試 randomString command 的執行
func TestRandomStringExecute(t *testing.T) {
	// 將標準輸出重定向到一個緩衝區
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 設定 flags 並執行 command
	randomStringLength = 12
	randomStringCount = 3
	randomStringCmd.Run(nil, nil)

	// 恢復標準輸出
	w.Close()
	os.Stdout = old

	// 從緩衝區讀取輸出
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// 執行斷言
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 3, "應產生正確數量的字串")
	for _, line := range lines {
		assert.Equal(t, 12, len(line), "每個產生的字串應具有正確的長度")
	}
}
