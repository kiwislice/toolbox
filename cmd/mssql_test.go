// 該檔案為 mssql command 的單元測試
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// mockColorPrinter 是 colorPrinter 介面的 mock 實作
type mockColorPrinter struct {
	buffer *bytes.Buffer
}

func (p *mockColorPrinter) Cyan(format string, a ...interface{}) {
	fmt.Fprintf(p.buffer, format, a...)
}
func (p *mockColorPrinter) Green(format string, a ...interface{}) {
	fmt.Fprintf(p.buffer, format, a...)
}
func (p *mockColorPrinter) Red(format string, a ...interface{}) {
	fmt.Fprintf(p.buffer, format, a...)
}

// resetMssqlCmdFlags 將 mssqlCmd 的所有旗標重設為預設值
func resetMssqlCmdFlags() {
	// 遍歷 mssqlCmd 的所有旗標
	mssqlCmd.Flags().VisitAll(func(f *pflag.Flag) {
		// 將旗標的值設為其預設值
		f.Value.Set(f.DefValue)
		// 將旗標的 "Changed" 狀態設為 false
		f.Changed = false
	})
}

// TestMssqlCmd_RequiredFlags 測試 mssql command 的必要參數
func TestMssqlCmd_RequiredFlags(t *testing.T) {
	// 準備一個 buffer 來捕捉 cobra 的輸出
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	// 測試案例
	testCases := []struct {
		name          string   // 測試案例名稱
		args          []string // 傳入的參數
		expectedError string   // 預期的錯誤訊息
	}{
		{
			name:          "缺少 user 參數",
			args:          []string{"mssql", "--password=pass", "--database=db"},
			expectedError: `required flag(s) "user" not set`,
		},
		{
			name:          "缺少 password 參數",
			args:          []string{"mssql", "--user=user", "--database=db"},
			expectedError: `required flag(s) "password" not set`,
		},
		{
			name:          "缺少 database 參數",
			args:          []string{"mssql", "--user=user", "--password=pass"},
			expectedError: `required flag(s) "database" not set`,
		},
	}

	// 遍歷所有測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 在每次測試前重設旗標
			resetMssqlCmdFlags()
			// 重設 buffer
			output.Reset()
			// 設定 rootCmd 的參數
			rootCmd.SetArgs(tc.args)
			// 執行 rootCmd，預期會回傳錯誤
			err := rootCmd.Execute()
			assert.Error(t, err)
			// 檢查錯誤訊息是否符合預期
			assert.True(t, strings.Contains(output.String(), tc.expectedError), "Expected error containing %q, but got %q", tc.expectedError, output.String())
		})
	}
}

// TestMssqlCmd_FlagsParsing 測試 mssql command 的參數解析
func TestMssqlCmd_FlagsParsing(t *testing.T) {
	// 備份原始的 printer 和 getSchema 函式
	originalPrinter := printer
	originalGetSchema := getSchema
	// 在函式結束時還原
	defer func() {
		printer = originalPrinter
		getSchema = originalGetSchema
	}()

	// 建立一個 buffer 來捕捉 mock printer 的輸出
	var mockPrinterBuffer bytes.Buffer
	// 建立一個 mock printer，並將其設定為全域 printer
	printer = &mockColorPrinter{buffer: &mockPrinterBuffer}

	// 建立一個假的 getSchema 函式，它會回傳一個空的 Schema 物件
	getSchema = func(info *tools.MssqlInfo) (*tools.Schema, error) {
		return &tools.Schema{}, nil
	}

	// 在測試前重設旗標
	resetMssqlCmdFlags()

	// 設定 rootCmd 的參數並執行
	rootCmd.SetArgs([]string{"mssql", "--host=testhost", "--port=1234", "--user=testuser", "--password=testpass", "--database=testdb", "--output=test_schema.md"})
	err := rootCmd.Execute()
	assert.NoError(t, err)

	output := mockPrinterBuffer.String()

	// 驗證輸出內容
	assert.Contains(t, output, "成功產生資料庫規格文件: test_schema.md")

	// 刪除測試產生的檔案
	os.Remove("test_schema.md")
}
