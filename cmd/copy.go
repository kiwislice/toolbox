package cmd

import (
	"errors"
	"path/filepath"

	"github.com/fatih/color"
	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy <source> <destination>",
	Short: "Copy files or folders",
	Long:  `Copy files or folders from a source to a destination.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dest := args[1]

		color.Cyan("開始複製：從 %s 到 %s", src, dest)
		defer color.Cyan("結束複製：從 %s 到 %s", src, dest)

		srcExist, srcInfo := tools.IsExist(src)
		if !srcExist {
			color.Red("來源不存在")
			return
		}
		destExist, destInfo := tools.IsExist(dest)

		srcType := tools.FileType(srcInfo)
		// 如果目標不存在則將目標當作跟來源同類
		destType := srcType
		if destExist {
			destType = tools.FileType(destInfo)
		}

		if srcType == tools.FILETYPE_UNKNOW || destType == tools.FILETYPE_UNKNOW {
			color.Red("來源或目的地類型不明")
			return
		}

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
		if err := function(src, dest); err != nil {
			tools.Errorf("%s", err)
		}
	},
}

func f2f(src, dest string) error {
	return tools.CopyFile(src, dest)
}

func f2d(src, dest string) error {
	filename := filepath.Base(src)
	dest = filepath.Join(dest, filename)
	return f2f(src, dest)
}

func d2f(src, dest string) error {
	return errors.New("資料夾無法複製到檔案")
}

func d2d(src, dest string) error {
	return tools.CopyDir(src, dest)
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
