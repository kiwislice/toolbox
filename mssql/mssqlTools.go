package mssql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	tools "github.com/kiwislice/toolbox/tools"
	"golang.org/x/exp/slices"
)

type CustomDB struct {
	*sql.DB
	dbname string
}

func NewDb(ip, port, acc, pw, dbname string) (db *CustomDB, err error) {
	// fmt.Println("connectToDb start")
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", ip, acc, pw, port, dbname)
	fmt.Println("connectToDb " + connString)

	sqldb, err := sql.Open("mssql", connString)
	if err != nil {
		tools.Errorf("sql.Open失敗: " + err.Error())
		return
	}
	return &CustomDB{DB: sqldb, dbname: dbname}, nil
}

func SelectAll(db *CustomDB, tablename string) (list []map[string]any, err error) {
	tsql := fmt.Sprintf("SELECT * FROM %s;", tablename)
	rows, err := db.Query(tsql)
	if err != nil {
		tools.Errorf("db.Query失敗: " + err.Error())
		return
	}
	defer rows.Close()
	return rows2MapArray(rows)
}

func rows2MapArray(rows *sql.Rows) (list []map[string]any, err error) {
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]any, columnLength) //临时存储每行数据
	for index := range cache {         //为每一列初始化一个指针
		var a any
		cache[index] = &a
	}

	for rows.Next() {
		err = rows.Scan(cache...)
		if err != nil {
			tools.Errorf("rows.Scan失敗: " + err.Error())
			break
		}
		item := make(map[string]any)
		for i, data := range cache {
			item[columns[i]] = fmt.Sprint(*data.(*any)) //取实际类型
		}
		list = append(list, item)
	}
	return list, err
}

type DatabaseSchema struct {
	Name   string
	Tables []*TableSchema
	Views  []*ViewSchema
}

type TableSchema struct {
	Name        string
	Columns     []*ColumnSchema
	PrimaryKey  *TableConstraint
	OtherKeys   []*TableConstraint
	Description string
}

type ViewSchema struct {
	Name        string
	Columns     []*ColumnSchema
	Description string
}

type ColumnSchema struct {
	TableName    string
	ColumnName   string
	Ordinal      int
	DefaultValue string
	Nullable     bool
	Datatype     string
	Description  string
	IsPrimaryKey bool
}

func (x *DatabaseSchema) ToJson() (string, error) {
	bar, err := json.Marshal(x)
	if err != nil {
		return "", err
	}
	return string(bar), nil
}

func GetDatabaseSchema(db *CustomDB) (databaseSchema *DatabaseSchema, err error) {
	sqlTemplate := `
SELECT tbl.TABLE_TYPE AS TABLE_TYPE
	,col.*
	,prop.value     AS [COLUMN_DESCRIPTION]
FROM %s.INFORMATION_SCHEMA.TABLES AS tbl
INNER JOIN INFORMATION_SCHEMA.COLUMNS AS col
ON col.TABLE_NAME = tbl.TABLE_NAME
INNER JOIN sys.columns AS sc
ON sc.object_id = object_id(tbl.table_schema + '.' + tbl.table_name) AND sc.NAME = col.COLUMN_NAME
LEFT JOIN sys.extended_properties prop
ON prop.major_id = sc.object_id AND prop.minor_id = sc.column_id AND prop.NAME = 'MS_Description'
WHERE tbl.TABLE_NAME <> 'sysdiagrams'
ORDER BY TABLE_NAME, ORDINAL_POSITION;`

	sql := fmt.Sprintf(sqlTemplate, db.dbname)
	rows, err := db.Query(sql)
	if err != nil {
		tools.Errorf("db.Query失敗: " + err.Error())
		return
	}
	defer rows.Close()

	list, err := rows2MapArray(rows)
	if err != nil {
		tools.Errorf("讀取rows失敗: " + err.Error())
		return
	}

	databaseSchema = new(DatabaseSchema)
	databaseSchema.Name = db.dbname
	tables := make(map[string]*TableSchema)
	views := make(map[string]*ViewSchema)
	for i := range list {
		// 讀取ColumnSchema
		cs := new(ColumnSchema)
		cs.TableName = list[i]["TABLE_NAME"].(string)
		cs.ColumnName = list[i]["COLUMN_NAME"].(string)
		cs.Ordinal, err = strconv.Atoi(list[i]["ORDINAL_POSITION"].(string))
		if err != nil {
			tools.Errorf("讀取cs.Ordinal失敗: " + err.Error())
		}
		if list[i]["COLUMN_DEFAULT"] != "<nil>" {
			cs.DefaultValue = list[i]["COLUMN_DEFAULT"].(string)
		}
		cs.Nullable = list[i]["IS_NULLABLE"].(string) == "YES"
		cs.Datatype = list[i]["DATA_TYPE"].(string)
		if strings.HasSuffix(cs.Datatype, "char") {
			length := list[i]["CHARACTER_MAXIMUM_LENGTH"].(string)
			if length == "-1" {
				length = "MAX"
			}
			cs.Datatype += fmt.Sprintf("(%s)", length)
		}
		if list[i]["COLUMN_DESCRIPTION"] != "<nil>" {
			cs.Description = list[i]["COLUMN_DESCRIPTION"].(string)
		}

		switch list[i]["TABLE_TYPE"] {
		case "BASE TABLE":
			ts, exist := tables[cs.TableName]
			if !exist {
				ts = new(TableSchema)
				ts.Name = cs.TableName
				tables[cs.TableName] = ts
				databaseSchema.Tables = append(databaseSchema.Tables, ts)
			}
			ts.Columns = append(ts.Columns, cs)
		case "VIEW":
			vs, exist := views[cs.TableName]
			if !exist {
				vs = new(ViewSchema)
				vs.Name = cs.TableName
				views[cs.TableName] = vs
				databaseSchema.Views = append(databaseSchema.Views, vs)
			}
			vs.Columns = append(vs.Columns, cs)
		}
	}

	// 補上條件約束資料
	constraints, err := getTableConstraints(db)
	if err != nil {
		tools.Errorf("讀取條件約束失敗: " + err.Error())
		return
	}

	for _, tableSchema := range databaseSchema.Tables {
		// show := tableSchema.Name == "act"
		for _, tableConstraint := range constraints[tableSchema.Name] {
			if tableConstraint.ConstraintType == "PRIMARY KEY" {
				tableSchema.PrimaryKey = tableConstraint
				for _, columnSchema := range tableSchema.Columns {
					columnSchema.IsPrimaryKey = slices.Contains(tableConstraint.Columns, columnSchema.ColumnName)
				}
			} else {
				tableSchema.OtherKeys = append(tableSchema.OtherKeys, tableConstraint)
			}
		}
		// if show {
		// 	jsonStr, err := json.Marshal(tableSchema)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	fmt.Println(string(jsonStr))
		// }
	}

	// 補上資料表&View的描述
	descriptions, err := getTableAndViewDescriptions(db)
	if err != nil {
		tools.Errorf("讀取資料表&View的描述失敗: " + err.Error())
		return
	}
	for _, tableSchema := range databaseSchema.Tables {
		tableSchema.Description = descriptions[tableSchema.Name]
	}
	for _, viewSchema := range databaseSchema.Views {
		viewSchema.Description = descriptions[viewSchema.Name]
	}

	return
}

// 資料表約束資料物件
type TableConstraint struct {
	TableName      string
	ConstraintName string
	ConstraintType string
	Columns        []string
}

// 取得所有資料表約束
//
// 回傳 map[資料表名稱][]TableConstraint
func getTableConstraints(db *CustomDB) (constraints map[string][]*TableConstraint, err error) {
	sql := `
SELECT KU.CONSTRAINT_NAME
	,TC.CONSTRAINT_TYPE
	,KU.table_name AS TABLENAME
	,column_name   AS KEYCOLUMN
FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC
INNER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS KU
ON TC.CONSTRAINT_NAME = KU.CONSTRAINT_NAME
ORDER BY KU.TABLE_NAME , KU.ORDINAL_POSITION;`

	rows, err := db.Query(sql)
	if err != nil {
		tools.Errorf("db.Query失敗: " + err.Error())
		return
	}
	defer rows.Close()

	list, err := rows2MapArray(rows)
	if err != nil {
		tools.Errorf("讀取rows失敗: " + err.Error())
		return
	}

	constraints = make(map[string][]*TableConstraint)
	var lastTC *TableConstraint
	for i := range list {
		tableName := list[i]["TABLENAME"].(string)
		constraintName := list[i]["CONSTRAINT_NAME"].(string)
		if lastTC == nil || lastTC.TableName != tableName || lastTC.ConstraintName != constraintName {
			lastTC = new(TableConstraint)
			lastTC.TableName = tableName
			lastTC.ConstraintName = constraintName
			lastTC.ConstraintType = list[i]["CONSTRAINT_TYPE"].(string)
			lastTC.Columns = make([]string, 0)

			constraints[tableName] = append(constraints[tableName], lastTC)
		}
		lastTC.Columns = append(lastTC.Columns, list[i]["KEYCOLUMN"].(string))
	}
	return
}

// 取得所有資料表&View描述(Description)
//
// 回傳 map[資料表orView名稱]描述
func getTableAndViewDescriptions(db *CustomDB) (m map[string]string, err error) {
	sql := `
SELECT  tvs.name   AS [TVNAME]
	,prop.value AS [DESCRIPTION]
FROM
(
 SELECT  name
		,object_id
 FROM sys.tables AS ts
 UNION
 SELECT  name
		,object_id
 FROM sys.views AS vs
) AS tvs
INNER JOIN sys.extended_properties prop
ON prop.major_id = tvs.object_id AND prop.minor_id = 0 AND prop.NAME = 'MS_Description';`

	rows, err := db.Query(sql)
	if err != nil {
		tools.Errorf("db.Query失敗: " + err.Error())
		return
	}
	defer rows.Close()

	list, err := rows2MapArray(rows)
	if err != nil {
		tools.Errorf("讀取rows失敗: " + err.Error())
		return
	}

	m = make(map[string]string)
	for i := range list {
		tableName := list[i]["TVNAME"].(string)
		description := list[i]["DESCRIPTION"].(string)
		m[tableName] = description
	}
	return
}
