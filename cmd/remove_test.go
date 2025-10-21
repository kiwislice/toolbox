// 該檔案為 remove command 的單元測試
package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRemoveFile 測試移除檔案
func TestRemoveFile(t *testing.T) {
	// 建立一個暫存檔案
	file, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	file.Close()

	// 執行 remove command
	removeCmd.Run(nil, []string{file.Name()})

	// 驗證檔案是否已被移除
	_, err = os.Stat(file.Name())
	assert.True(t, os.IsNotExist(err))
}

// TestRemoveDir 測試移除資料夾
func TestRemoveDir(t *testing.T) {
	// 建立一個暫存資料夾
	dir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	// 執行 remove command
	removeCmd.Run(nil, []string{dir})

	// 驗證資料夾是否已被移除
	_, err = os.Stat(dir)
	assert.True(t, os.IsNotExist(err))
}

// TestRemoveNotExist 測試移除不存在的檔案或資料夾 (預期不會發生錯誤)
func TestRemoveNotExist(t *testing.T) {
	removeCmd.Run(nil, []string{"not-exist"})
}
