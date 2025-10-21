// 整個 toolbox 工具的進入點
package main

import (
	"github.com/kiwislice/toolbox/cmd"
)

// main 函式，程式的進入點
func main() {
	// 執行 cmd 套件中的 Execute 函式
	cmd.Execute()
}
