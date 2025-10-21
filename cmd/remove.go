// 該檔案定義了 remove command
package cmd

import (
	"os"

	"github.com/fatih/color"
	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

// removeCmd 是 `toolbox remove` command 的定義
var removeCmd = &cobra.Command{
	Use:   "remove <target>",
	Short: "移除檔案或資料夾",
	Long:  `移除指定的檔案或資料夾`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 取得使用者輸入的目標路徑
		target := args[0]

		color.Cyan("開始刪除：%s", target)
		defer color.Cyan("結束刪除：%s", target)

		// 執行刪除
		err := os.RemoveAll(target)
		if err != nil {
			tools.Errorf("%s", err)
		}
	},
}

func init() {
	// 將 removeCmd 加入 rootCmd 中
	rootCmd.AddCommand(removeCmd)
}
