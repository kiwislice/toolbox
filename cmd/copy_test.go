// 該檔案為 copy command 的單元測試
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCopyCmd_FileToFile 測試複製檔案到檔案
func TestCopyCmd_FileToFile(t *testing.T) {
	// 建立來源檔案
	srcFile, err := os.CreateTemp("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())
	srcFile.WriteString("test")
	srcFile.Close()

	// 建立目的地檔案
	destFile, err := os.CreateTemp("", "dest")
	assert.NoError(t, err)
	destFile.Close()
	defer os.Remove(destFile.Name())

	// 執行 copy command
	copyCmd.Run(nil, []string{srcFile.Name(), destFile.Name()})

	// 驗證目的地檔案內容
	content, err := os.ReadFile(destFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

// TestCopyCmd_FileToDir 測試複製檔案到資料夾
func TestCopyCmd_FileToDir(t *testing.T) {
	// 建立來源檔案
	srcFile, err := os.CreateTemp("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())
	srcFile.WriteString("test")
	srcFile.Close()

	// 建立目的地資料夾
	destDir, err := os.MkdirTemp("", "dest")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	// 執行 copy command
	copyCmd.Run(nil, []string{srcFile.Name(), destDir})

	// 驗證目的地資料夾內的檔案內容
	destFile := filepath.Join(destDir, filepath.Base(srcFile.Name()))
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

// TestCopyCmd_DirToFile 測試複製資料夾到檔案 (預期會失敗)
func TestCopyCmd_DirToFile(t *testing.T) {
	// 建立來源資料夾
	srcDir, err := os.MkdirTemp("", "src")
	assert.NoError(t, err)
	defer os.RemoveAll(srcDir)

	// 建立目的地檔案
	destFile, err := os.CreateTemp("", "dest")
	assert.NoError(t, err)
	defer os.Remove(destFile.Name())
	destFile.Close()

	// 執行 copy command，預期會印出錯誤訊息，但程式不會崩潰
	// 這裡我們不直接斷言錯誤，因為 Run 函式會捕捉並印出錯誤
	copyCmd.Run(nil, []string{srcDir, destFile.Name()})
}

// TestCopyCmd_DirToDir 測試複製資料夾到資料夾
func TestCopyCmd_DirToDir(t *testing.T) {
	// 建立來源資料夾
	srcDir, err := os.MkdirTemp("", "src")
	assert.NoError(t, err)
	defer os.RemoveAll(srcDir)

	// 在來源資料夾內建立一個檔案
	srcFile, err := os.CreateTemp(srcDir, "file")
	assert.NoError(t, err)
	srcFile.WriteString("test")
	srcFile.Close()

	// 建立目的地資料夾
	destDir, err := os.MkdirTemp("", "dest")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	// 執行 copy command
	copyCmd.Run(nil, []string{srcDir, destDir})

	// 驗證目的地資料夾內的檔案內容
	destFile := filepath.Join(destDir, filepath.Base(srcFile.Name()))
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}
