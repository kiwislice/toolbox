// 該檔案為 fileTools 的單元測試
package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_IsExist 測試 IsExist 函式
func Test_IsExist(t *testing.T) {
	fmt.Println("test start.")
	assert := assert.New(t)

	// 取得目前的工作目錄
	workDir, err := os.Getwd()
	if err != nil {
		assert.Fail(fmt.Sprint(err))
		return
	}

	// 測試函式，用來驗證 IsExist 的結果
	test := func(path string) {
		expected, fileInfo := IsExist(path)

		if expected {
			if fileInfo.IsDir() {
				assert.DirExists(path)
			} else {
				assert.FileExists(path)
			}
		} else {
			assert.NoFileExists(path)
			assert.NoDirExists(path)
		}
	}

	// 測試一個不存在的檔案
	path := filepath.Join(workDir, "testfile")
	test(path)

	// 建立一個檔案，再測試一次
	f, err := os.Create(path)
	if err != nil {
		assert.Fail(fmt.Sprint(err))
		return
	}
	defer os.RemoveAll(path)
	f.Close()
	test(path)

	fmt.Println("test finished.")
	// 這個測試案例會刻意失敗，這是已知的狀況
	assert.FailNow("a")
}
