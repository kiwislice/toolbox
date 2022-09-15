package tools

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	FILETYPE_FILE = iota
	FILETYPE_DIR
	FILETYPE_UNKNOW
)

type FileInfo = fs.FileInfo

func GetFileInfo(path string) (FileInfo, error) {
	return os.Stat(path)
}

func IsExist(path string) (bool, FileInfo) {
	if fileInfo, err := GetFileInfo(path); err == nil {
		// path exists
		return true, fileInfo
	} else if errors.Is(err, os.ErrNotExist) {
		// path does *not* exist
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, nil
	}
}

func FileType(fileInfo FileInfo) int {
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		return FILETYPE_DIR
	case mode.IsRegular():
		return FILETYPE_FILE
	}
	return FILETYPE_UNKNOW
}

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

func CopyDir(src, dst string) error {
	e := os.MkdirAll(dst, os.ModePerm)
	if e != nil {
		return e
	}
	pb := NewLoadingText()
	pb.Start()

	fn := func(path string, d fs.DirEntry, err error) error {
		pb.SetText(path)

		Debugf("walk to: %s", path)
		relPath, e := filepath.Rel(src, path)
		if e != nil {
			return e
		}
		targetPath := filepath.Join(dst, relPath)
		Debugf("target: %s", path)

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
