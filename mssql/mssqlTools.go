package mssql

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	tools "github.com/kiwislice/toolbox/tools"
)

func NewDb(ip, port, acc, pw, dbname string) (db *sql.DB, err error) {
	// fmt.Println("connectToDb start")
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", ip, acc, pw, port, dbname)
	fmt.Println("connectToDb " + connString)
	return sql.Open("mssql", connString)
}

func SelectAll(db *sql.DB, tablename string) (list []map[string]any, err error) {
	tsql := fmt.Sprintf("SELECT * FROM %s;", tablename)
	rows, err := db.Query(tsql)
	if err != nil {
		tools.Errorf("db.Query失敗: " + err.Error())
		return
	}
	defer rows.Close()

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
