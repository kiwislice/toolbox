package clear

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClearDir(t *testing.T) {
	// 建立一個暫存目錄
	dir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(dir) // 確保測試結束後刪除目錄

	// 在目錄中建立一個檔案
	tmpfn := filepath.Join(dir, "tmpfile")
	err = os.WriteFile(tmpfn, []byte("test data"), 0666)
	assert.NoError(t, err)

	// 執行 clear 指令
	cmd, err := newClearCmd([]string{dir})
	assert.NoError(t, err)
	cmd.Run()

	// 檢查目錄是否存在
	_, err = os.Stat(dir)
	assert.NoError(t, err, "Directory should still exist")

	// 檢查目錄是否為空
	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)
	assert.Empty(t, entries, "Directory should be empty")
}

func TestClearNonExistentDir(t *testing.T) {
	// 對一個不存在的目錄執行 clear 指令
	nonExistentDir := "non-existent-dir-for-test"
	cmd, err := newClearCmd([]string{nonExistentDir})
	assert.NoError(t, err)
	cmd.Run()

	// 檢查該目錄是否仍然不存在
	_, err = os.Stat(nonExistentDir)
	assert.True(t, os.IsNotExist(err), "Non-existent directory should not be created")
}

func TestClearFile(t *testing.T) {
	// 建立一個暫存檔案
	file, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(file.Name()) // 確保測試結束後刪除檔案
	file.Close()

	// 對檔案執行 clear 指令
	cmd, err := newClearCmd([]string{file.Name()})
	assert.NoError(t, err)
	cmd.Run()

	// 檢查檔案是否仍然存在
	_, err = os.Stat(file.Name())
	assert.NoError(t, err, "File should not be deleted")
}
