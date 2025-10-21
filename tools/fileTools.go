// 該檔案提供檔案/資料夾操作相關的工具函式
package tools

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// 檔案類型的常數
const (
	FILETYPE_FILE   = iota // 檔案
	FILETYPE_DIR           // 資料夾
	FILETYPE_UNKNOW        // 未知
)

// FileInfo fs.FileInfo 的別名
type FileInfo = fs.FileInfo

// GetFileInfo 取得檔案/資料夾的資訊
func GetFileInfo(path string) (FileInfo, error) {
	return os.Stat(path)
}

// IsExist 檢查檔案/資料夾是否存在
func IsExist(path string) (bool, FileInfo) {
	if fileInfo, err := GetFileInfo(path); err == nil {
		// 檔案/資料夾存在
		return true, fileInfo
	} else if errors.Is(err, os.ErrNotExist) {
		// 檔案/資料夾不存在
		return false, nil
	} else {
		// 其他錯誤
		return false, nil
	}
}

// FileType 判斷檔案類型 (檔案或資料夾)
func FileType(fileInfo FileInfo) int {
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		return FILETYPE_DIR
	case mode.IsRegular():
		return FILETYPE_FILE
	}
	return FILETYPE_UNKNOW
}

// CopyFile 複製檔案
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

// CopyDir 複製資料夾
func CopyDir(src, dst string) error {
	// 建立目的地資料夾
	e := os.MkdirAll(dst, os.ModePerm)
	if e != nil {
		return e
	}
	// 顯示進度條
	pb := NewLoadingText()
	pb.Start()

	// 遍歷來源資料夾內的所有檔案/資料夾
	fn := func(path string, d fs.DirEntry, err error) error {
		pb.SetText(path)

		Debugf("walk to: %s", path)
		// 取得相對路徑
		relPath, e := filepath.Rel(src, path)
		if e != nil {
			return e
		}
		// 組成目標路徑
		targetPath := filepath.Join(dst, relPath)
		Debugf("target: %s", path)

		// 如果是資料夾就建立，如果是檔案就複製
		if d.IsDir() {
			e = os.MkdirAll(targetPath, os.ModePerm)
		} else {
			e = CopyFile(path, targetPath)
		}

		return e
	}
	e = filepath.WalkDir(src, fn)
	pb.Finish()
	return e
}

// WriteTextFile 新增或覆蓋文字檔
func WriteTextFile(text, path string) (err error) {
	// 建立檔案
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("os.Create失敗:", err.Error())
		return
	}
	defer f.Close()

	// 寫入文字
	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println("f.WriteString失敗:", err.Error())
	}
	return
}
