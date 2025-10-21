// 該檔案定義了 clear command
package cmd

import (
	"os"

	"github.com/fatih/color"
	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

// clearCmd 是 `toolbox clear` command 的定義
var clearCmd = &cobra.Command{
	Use:   "clear <target>",
	Short: "清空一個資料夾",
	Long:  `清空一個資料夾`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 取得使用者輸入的目標路徑
		target := args[0]

		color.Cyan("開始清空：%s", target)
		defer color.Cyan("結束清空：%s", target)

		// 檢查目標是否存在且為資料夾
		exist, fileInfo := tools.IsExist(target)
		if !exist || !fileInfo.IsDir() {
			tools.Warnf("目標不是資料夾：%s", target)
			return
		}

		// 移除目標資料夾內所有檔案
		err := os.RemoveAll(target)
		if err != nil {
			tools.Errorf("%s", err)
		}

		// 重新建立目標資料夾
		err = os.Mkdir(target, os.ModePerm)
		if err != nil {
			tools.Errorf("%s", err)
		}
	},
}

// init 函式會在程式啟動時自動執行
func init() {
	// 將 clearCmd 加入 rootCmd 中
	rootCmd.AddCommand(clearCmd)
}
