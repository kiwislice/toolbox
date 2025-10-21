// 該檔案定義了 copy command
package cmd

import (
	"errors"
	"path/filepath"

	"github.com/fatih/color"
	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

// copyCmd 是 `toolbox copy` command 的定義
var copyCmd = &cobra.Command{
	Use:   "copy <source> <destination>",
	Short: "複製檔案或資料夾",
	Long:  `從來源複製檔案或資料夾到目的地`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// 取得使用者輸入的來源與目的地路徑
		src := args[0]
		dest := args[1]

		color.Cyan("開始複製：從 %s 到 %s", src, dest)
		defer color.Cyan("結束複製：從 %s 到 %s", src, dest)

		// 檢查來源是否存在
		srcExist, srcInfo := tools.IsExist(src)
		if !srcExist {
			color.Red("來源不存在")
			return
		}
		// 檢查目的地是否存在
		destExist, destInfo := tools.IsExist(dest)

		// 取得來源的檔案類型 (檔案或資料夾)
		srcType := tools.FileType(srcInfo)
		// 如果目標不存在則將目標當作跟來源同類
		destType := srcType
		if destExist {
			destType = tools.FileType(destInfo)
		}

		// 如果來源或目的地的類型不明，則無法處理
		if srcType == tools.FILETYPE_UNKNOW || destType == tools.FILETYPE_UNKNOW {
			color.Red("來源或目的地類型不明")
			return
		}

		// 根據來源與目的地的類型，決定要執行的複製函式
		var function func(src, dest string) error
		if srcType == tools.FILETYPE_FILE {
			if destType == tools.FILETYPE_FILE {
				function = f2f
			} else {
				function = f2d
			}
		} else {
			if destType == tools.FILETYPE_FILE {
				function = d2f
			} else {
				function = d2d
			}
		}

		// 執行複製函式
		if err := function(src, dest); err != nil {
			tools.Errorf("%s", err)
		}
	},
}

// f2f (file to file) 複製檔案到檔案
func f2f(src, dest string) error {
	return tools.CopyFile(src, dest)
}

// f2d (file to directory) 複製檔案到資料夾
func f2d(src, dest string) error {
	filename := filepath.Base(src)
	dest = filepath.Join(dest, filename)
	return f2f(src, dest)
}

// d2f (directory to file) 複製資料夾到檔案 (不支援)
func d2f(src, dest string) error {
	return errors.New("資料夾無法複製到檔案")
}

// d2d (directory to directory) 複製資料夾到資料夾
func d2d(src, dest string) error {
	return tools.CopyDir(src, dest)
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
