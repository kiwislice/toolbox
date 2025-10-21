// cmd 套件包含了所有 cobra command 的定義
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 是整個程式的主要 command
var rootCmd = &cobra.Command{
	Use:   "toolbox",
	Short: "一個多功能工具箱",
	Long: `一個多功能工具箱，包含多種實用工具。
例如：
- 清除螢幕
- 複製檔案
- 產生隨機字串
- 移除檔案`,
}

// Execute 函式會執行 rootCmd
func Execute() {
	// 執行 rootCmd，如果發生錯誤則會印出錯誤訊息並結束程式
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
