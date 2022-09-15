package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_ShowPath(t *testing.T) {
// 	fmt.Println("Test_ShowPath start.")
// 	assert := assert.New(t)

// 	os.Remove()

// 	f, err := os.CreateTemp("tempfile")
// 	if err != nil {
// 		assert.Fail(fmt.Sprint(err))
// 		return
// 	}
// 	defer f.Close()

// 	wd, err := os.Getwd()
// 	if err != nil {
// 		assert.Fail(fmt.Sprint(err))
// 		return
// 	}

// 	fmt.Println("test finished." + wd)
// 	assert.FailNow("Test_ShowPath")
// }

func Test_IsExist(t *testing.T) {
	fmt.Println("test start.")
	assert := assert.New(t)

	workDir, err := os.Getwd()
	if err != nil {
		assert.Fail(fmt.Sprint(err))
		return
	}

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

	path := filepath.Join(workDir, "testfile")
	test(path)

	f, err := os.Create(path)
	if err != nil {
		assert.Fail(fmt.Sprint(err))
		return
	}
	defer os.RemoveAll(path)
	f.Close()
	test(path)

	fmt.Println("test finished.")
	assert.FailNow("a")
}
