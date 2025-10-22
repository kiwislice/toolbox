// 該檔案定義了 mssql command
package cmd

import (
	"os"

	"github.com/fatih/color"
	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

// colorPrinter 是一個介面，定義了印出訊息的行為
type colorPrinter interface {
	Cyan(format string, a ...interface{})
	Green(format string, a ...interface{})
	Red(format string, a ...interface{})
}

// defaultColorPrinter 是 colorPrinter 介面的預設實作
type defaultColorPrinter struct{}

func (p *defaultColorPrinter) Cyan(format string, a ...interface{}) {
	color.Cyan(format, a...)
}
func (p *defaultColorPrinter) Green(format string, a ...interface{}) {
	color.Green(format, a...)
}
func (p *defaultColorPrinter) Red(format string, a ...interface{}) {
	color.Red(format, a...)
}

// printer 是我們在程式中使用的 colorPrinter 實例
var printer colorPrinter = &defaultColorPrinter{}

// getSchema 是一個變數，其預設值為 (*tools.MssqlInfo).GetSchema 函式
// 這樣我們就可以在測試中覆寫它，以注入一個 mock 的 GetSchema 函式
var getSchema = (*tools.MssqlInfo).GetSchema

// mssqlCmd 是 `toolbox mssql` command 的定義
var mssqlCmd = &cobra.Command{
	Use:   "mssql",
	Short: "產生 MSSQL 資料庫規格文件",
	Long:  `連線到指定的 MSSQL 資料庫，並產生一份完整的資料庫規格文件。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 從 command line flags 讀取資料庫連線資訊
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		output, _ := cmd.Flags().GetString("output")

		// 建立 MssqlInfo 物件
		info := &tools.MssqlInfo{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			Database: database,
		}

		printer.Cyan("開始產生資料庫規格文件...")

		// 獲取資料庫結構
		schema, err := getSchema(info)
		if err != nil {
			printer.Red("獲取資料庫結構時發生錯誤: %v", err)
			return
		}

		// 產生 Markdown 文件
		markdown := schema.GenerateMarkdown()

		// 將 Markdown 文件寫入檔案
		err = os.WriteFile(output, []byte(markdown), 0644)
		if err != nil {
			printer.Red("寫入檔案時發生錯誤: %v", err)
			return
		}

		printer.Green("成功產生資料庫規格文件: %s", output)
	},
}

func init() {
	// 將 mssqlCmd 加入到 rootCmd
	rootCmd.AddCommand(mssqlCmd)

	// 為 mssqlCmd 定義 command line flags
	mssqlCmd.Flags().StringP("host", "H", "localhost", "資料庫主機")
	mssqlCmd.Flags().IntP("port", "P", 1433, "資料庫連接埠")
	mssqlCmd.Flags().StringP("user", "u", "", "使用者名稱")
	mssqlCmd.Flags().StringP("password", "p", "", "密碼")
	mssqlCmd.Flags().StringP("database", "d", "", "資料庫名稱")
	mssqlCmd.Flags().StringP("output", "o", "database_schema.md", "輸出檔案路徑")


	// 將 'user', 'password', 'database' 設為必填參數
	mssqlCmd.MarkFlagRequired("user")
	mssqlCmd.MarkFlagRequired("password")
	mssqlCmd.MarkFlagRequired("database")
}
