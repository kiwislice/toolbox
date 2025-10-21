package cmd

import (
	"os"

	"github.com/fatih/color"
	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear <target>",
	Short: "Clear a folder",
	Long:  `Clear a folder.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		color.Cyan("開始清空：%s", target)
		defer color.Cyan("結束清空：%s", target)

		exist, fileInfo := tools.IsExist(target)
		if !exist || !fileInfo.IsDir() {
			tools.Warnf("目標不是資料夾：%s", target)
			return
		}

		err := os.RemoveAll(target)
		if err != nil {
			tools.Errorf("%s", err)
		}
		err = os.Mkdir(target, os.ModePerm)
		if err != nil {
			tools.Errorf("%s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
