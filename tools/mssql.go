// 該檔案提供了與 MSSQL 資料庫互動的工具函式
package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"strings"
)

// sqlOpen 是一個變數，其預設值為 sql.Open 函式
// 這樣我們就可以在測試中覆寫它，以注入一個 mock 的資料庫連線
var sqlOpen = sql.Open

// MssqlInfo 包含了連線 MSSQL 資料庫所需的所有資訊
type MssqlInfo struct {
	Host     string // 資料庫主機
	Port     int    // 資料庫連接埠
	User     string // 使用者名稱
	Password string // 密碼
	Database string // 資料庫名稱
}

// Column 代表資料庫中的一個欄位
type Column struct {
	Name            string // 欄位名稱
	OrdinalPosition int    // 欄位順序
	IsNullable      string // 是否可為空
	DataType        string // 資料型態
	MaxLength       sql.NullInt64 // 最大長度
	Precision       sql.NullInt64 // 數值精度
	Scale           sql.NullInt64 // 數值小數位數
	Default         sql.NullString // 預設值
	IsPrimaryKey    bool   // 是否為主鍵
	Description     string // 欄位說明
}

// Table 代表資料庫中的一個資料表
type Table struct {
	Schema      string   // 資料表所屬的 Schema
	Name        string   // 資料表名稱
	Description string   // 資料表說明
	Columns     []Column // 資料表的所有欄位
}

// Schema 代表整個資料庫的結構
type Schema struct {
	Tables []Table // 資料庫中的所有資料表
}


// Connect 函式會根據提供的 MssqlInfo 建立一個資料庫連線
// 它會回傳一個 *sql.DB 物件，如果連線失敗則會回傳錯誤
func (info *MssqlInfo) Connect() (*sql.DB, error) {
	// 建立連線字串
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		info.Host, info.User, info.Password, info.Port, info.Database)

	// 使用我們定義的 sqlOpen 變數來開啟資料庫連線
	db, err := sqlOpen("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("開啟資料庫連線時發生錯誤: %w", err)
	}

	// 測試資料庫連線是否成功
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("無法連線到資料庫: %w", err)
	}

	// 回傳資料庫連線物件
	return db, nil
}

// GetSchema 函式會連線到資料庫，並獲取其完整的結構資訊
func (info *MssqlInfo) GetSchema() (*Schema, error) {
	// 建立資料庫連線
	db, err := info.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 建立一個 Schema 物件來儲存結果
	schema := &Schema{}

	// 獲取所有資料表
	tables, err := getTables(db)
	if err != nil {
		return nil, err
	}

	// 遍歷所有資料表，獲取其欄位資訊
	for i := range tables {
		columns, err := getColumns(db, tables[i].Schema, tables[i].Name)
		if err != nil {
			// 如果查詢某個資料表的欄位時發生錯誤，只記錄錯誤，不中斷整個流程
			log.Printf("查詢資料表 '%s.%s' 的欄位時發生錯誤: %v", tables[i].Schema, tables[i].Name, err)
			continue
		}
		tables[i].Columns = columns
	}
	schema.Tables = tables

	return schema, nil
}


// getTables 函式會查詢並回傳資料庫中所有的資料表
func getTables(db *sql.DB) ([]Table, error) {
	// SQL 查詢語法，用於獲取所有資料表及其說明
	query := `
		SELECT
			t.TABLE_SCHEMA,
			t.TABLE_NAME,
			p.value AS 'TableDescription'
		FROM
			INFORMATION_SCHEMA.TABLES t
		LEFT JOIN
			sys.extended_properties p ON p.major_id = OBJECT_ID(t.TABLE_SCHEMA + '.' + t.TABLE_NAME)
									 AND p.minor_id = 0
									 AND p.name = 'MS_Description'
		WHERE
			t.TABLE_TYPE = 'BASE TABLE'
		ORDER BY
			t.TABLE_SCHEMA, t.TABLE_NAME;
	`

	// 執行查詢
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查詢資料表時發生錯誤: %w", err)
	}
	defer rows.Close()

	// 建立一個 slice 來儲存資料表
	var tables []Table
	// 遍歷查詢結果
	for rows.Next() {
		var table Table
		var description sql.NullString // 說明可能為 NULL
		// 掃描結果到 table 物件
		if err := rows.Scan(&table.Schema, &table.Name, &description); err != nil {
			log.Printf("掃描資料表結果時發生錯誤: %v", err)
			continue
		}
		// 如果說明存在，則將其設定到 table 物件
		if description.Valid {
			table.Description = description.String
		}
		// 將 table 物件加入到 slice
		tables = append(tables, table)
	}

	return tables, nil
}

// getColumns 函式會查詢並回傳指定資料表的所有欄位資訊
func getColumns(db *sql.DB, tableSchema, tableName string) ([]Column, error) {
	// SQL 查詢語法，用於獲取指定資料表的欄位資訊
	query := `
		SELECT
			c.COLUMN_NAME,
			c.ORDINAL_POSITION,
			c.IS_NULLABLE,
			c.DATA_TYPE,
			c.CHARACTER_MAXIMUM_LENGTH,
			c.NUMERIC_PRECISION,
			c.NUMERIC_SCALE,
			c.COLUMN_DEFAULT,
			p.value AS 'ColumnDescription',
			CASE WHEN pk.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END AS 'IsPrimaryKey'
		FROM
			INFORMATION_SCHEMA.COLUMNS c
		LEFT JOIN
			sys.extended_properties p ON p.major_id = OBJECT_ID(c.TABLE_SCHEMA + '.' + c.TABLE_NAME)
									 AND p.minor_id = c.ORDINAL_POSITION
									 AND p.name = 'MS_Description'
		LEFT JOIN
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE pk ON pk.TABLE_SCHEMA = c.TABLE_SCHEMA
												   AND pk.TABLE_NAME = c.TABLE_NAME
												   AND pk.COLUMN_NAME = c.COLUMN_NAME
												   AND OBJECTPROPERTY(OBJECT_ID(pk.CONSTRAINT_SCHEMA + '.' + pk.CONSTRAINT_NAME), 'IsPrimaryKey') = 1
		WHERE
			c.TABLE_SCHEMA = @p1 AND c.TABLE_NAME = @p2
		ORDER BY
			c.ORDINAL_POSITION;
	`
	// 執行查詢
	rows, err := db.Query(query, tableSchema, tableName)
	if err != nil {
		return nil, fmt.Errorf("查詢欄位時發生錯誤: %w", err)
	}
	defer rows.Close()

	// 建立一個 slice 來儲存欄位
	var columns []Column
	// 遍歷查詢結果
	for rows.Next() {
		var col Column
		var description sql.NullString // 說明可能為 NULL
		// 掃描結果到 col 物件
		err := rows.Scan(
			&col.Name,
			&col.OrdinalPosition,
			&col.IsNullable,
			&col.DataType,
			&col.MaxLength,
			&col.Precision,
			&col.Scale,
			&col.Default,
			&description,
			&col.IsPrimaryKey,
		)
		if err != nil {
			log.Printf("掃描欄位結果時發生錯誤: %v", err)
			continue
		}
		// 如果說明存在，則將其設定到 col 物件
		if description.Valid {
			col.Description = description.String
		}
		// 將 col 物件加入到 slice
		columns = append(columns, col)
	}
	return columns, nil
}

// GenerateMarkdown 函式會將 Schema 物件轉換成 Markdown 格式的字串
func (s *Schema) GenerateMarkdown() string {
	var builder strings.Builder

	// 文件標題
	builder.WriteString("# 資料庫規格文件\n\n")

	// 遍歷所有資料表
	for _, table := range s.Tables {
		// 資料表名稱和說明
		builder.WriteString(fmt.Sprintf("## %s.%s\n", table.Schema, table.Name))
		if table.Description != "" {
			builder.WriteString(fmt.Sprintf("%s\n\n", table.Description))
		}

		// 欄位表格的標頭
		builder.WriteString("| 欄位名稱 | 資料型態 | 可否為空 | 主鍵 | 預設值 | 說明 |\n")
		builder.WriteString("|---|---|---|---|---|---|\n")

		// 遍歷所有欄位
		for _, col := range table.Columns {
			// 格式化資料型態和長度
			dataType := col.DataType
			if col.MaxLength.Valid {
				dataType = fmt.Sprintf("%s(%d)", dataType, col.MaxLength.Int64)
			} else if col.Precision.Valid && col.Scale.Valid {
				dataType = fmt.Sprintf("%s(%d, %d)", dataType, col.Precision.Int64, col.Scale.Int64)
			}

			// 格式化主鍵
			isPrimaryKey := ""
			if col.IsPrimaryKey {
				isPrimaryKey = "✓"
			}

			// 格式化預設值
			defaultValue := ""
			if col.Default.Valid {
				defaultValue = col.Default.String
			}

			// 將欄位資訊寫入表格
			builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				col.Name,
				dataType,
				col.IsNullable,
				isPrimaryKey,
				defaultValue,
				col.Description,
			))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
