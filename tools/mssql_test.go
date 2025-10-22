// 該檔案為 tools/mssql.go 的單元測試
package tools

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestMssqlInfo_Connect_Success 測試成功的資料庫連線
func TestMssqlInfo_Connect_Success(t *testing.T) {
	// 建立一個 sqlmock 的資料庫連線和 mock 物件，並啟用 Ping 監控
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	// 預期 mock 會收到一個 Ping() 的請求，並且不會回傳任何錯誤
	mock.ExpectPing().WillReturnError(nil)

	// 建立一個 MssqlInfo 物件
	info := &MssqlInfo{
		Host:     "localhost",
		Port:     1433,
		User:     "user",
		Password: "password",
		Database: "database",
	}

	// 建立一個假的連線函式，它會直接回傳 mock 的資料庫連線
	// 這樣我們就可以繞過實際的 sql.Open()，避免真的去連線資料庫
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	// 呼叫 Connect() 函式
	_, err = info.Connect()

	// 驗證結果
	assert.NoError(t, err)
	// 確保所有預期的請求都已經被滿足
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestMssqlInfo_Connect_Failure 測試失敗的資料庫連線
func TestMssqlInfo_Connect_Failure(t *testing.T) {
	// 備份原始的 sqlOpen 函式
	originalSqlOpen := sqlOpen
	// 在函式結束時還原
	defer func() { sqlOpen = originalSqlOpen }()

	// 建立一個假的連線函式，它會回傳一個錯誤
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("連線失敗")
	}

	// 建立一個 MssqlInfo 物件
	info := &MssqlInfo{
		Host:     "invalid-host",
		Port:     9999,
		User:     "invalid-user",
		Password: "invalid-password",
		Database: "invalid-database",
	}

	// 執行 Connect() 函式，預期會回傳一個錯誤
	_, err := info.Connect()

	// 驗證結果
	assert.Error(t, err)
}

// TestGetSchema 測試獲取資料庫結構的邏輯
func TestGetSchema(t *testing.T) {
	// 建立一個 sqlmock 的資料庫連線和 mock 物件
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 模擬 getTables 的查詢結果
	tablesRows := sqlmock.NewRows([]string{"TABLE_SCHEMA", "TABLE_NAME", "TableDescription"}).
		AddRow("dbo", "Users", "使用者資料表")
	mock.ExpectQuery("^SELECT\\s+t.TABLE_SCHEMA,\\s+t.TABLE_NAME,\\s+p.value AS 'TableDescription'").WillReturnRows(tablesRows)

	// 模擬 getColumns 的查詢結果
	columnsRows := sqlmock.NewRows([]string{"COLUMN_NAME", "ORDINAL_POSITION", "IS_NULLABLE", "DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE", "COLUMN_DEFAULT", "ColumnDescription", "IsPrimaryKey"}).
		AddRow("ID", 1, "NO", "int", nil, 10, 0, nil, "使用者 ID", true).
		AddRow("Name", 2, "YES", "varchar", 50, nil, nil, "('John Doe')", "使用者名稱", false)
	mock.ExpectQuery("^SELECT\\s+c.COLUMN_NAME,\\s+c.ORDINAL_POSITION,\\s+c.IS_NULLABLE,\\s+c.DATA_TYPE,").WithArgs("dbo", "Users").WillReturnRows(columnsRows)

	// 建立一個假的連線函式
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	// 建立一個 MssqlInfo 物件
	info := &MssqlInfo{
		Host:     "localhost",
		Port:     1433,
		User:     "user",
		Password: "password",
		Database: "database",
	}

	// 呼叫 GetSchema 函式
	schema, err := info.GetSchema()

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Len(t, schema.Tables, 1)
	assert.Equal(t, "dbo", schema.Tables[0].Schema)
	assert.Equal(t, "Users", schema.Tables[0].Name)
	assert.Equal(t, "使用者資料表", schema.Tables[0].Description)
	assert.Len(t, schema.Tables[0].Columns, 2)
	assert.Equal(t, "ID", schema.Tables[0].Columns[0].Name)
	assert.True(t, schema.Tables[0].Columns[0].IsPrimaryKey)
	assert.Equal(t, "Name", schema.Tables[0].Columns[1].Name)
	assert.False(t, schema.Tables[0].Columns[1].IsPrimaryKey)
}

// TestGenerateMarkdown 測試產生 Markdown 文件的邏輯
func TestGenerateMarkdown(t *testing.T) {
	// 建立一個假的 Schema 物件
	schema := &Schema{
		Tables: []Table{
			{
				Schema:      "dbo",
				Name:        "Users",
				Description: "使用者資料表",
				Columns: []Column{
					{Name: "ID", DataType: "int", IsNullable: "NO", IsPrimaryKey: true, Description: "使用者 ID"},
					{Name: "Name", DataType: "varchar", MaxLength: sql.NullInt64{Int64: 50, Valid: true}, IsNullable: "YES", IsPrimaryKey: false, Default: sql.NullString{String: "('John Doe')", Valid: true}, Description: "使用者名稱"},
				},
			},
		},
	}

	// 呼叫 GenerateMarkdown 函式
	markdown := schema.GenerateMarkdown()

	// 驗證結果
	assert.Contains(t, markdown, "# 資料庫規格文件")
	assert.Contains(t, markdown, "## dbo.Users")
	assert.Contains(t, markdown, "使用者資料表")
	assert.Contains(t, markdown, "| 欄位名稱 | 資料型態 | 可否為空 | 主鍵 | 預設值 | 說明 |")
	assert.Contains(t, markdown, "| ID | int | NO | ✓ |  | 使用者 ID |")
	assert.Contains(t, markdown, "| Name | varchar(50) | YES |  | ('John Doe') | 使用者名稱 |")
}
