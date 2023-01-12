package mssql

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	core "github.com/kiwislice/toolbox/core"
	tools "github.com/kiwislice/toolbox/tools"
)

/// 取得資料庫schema的指令物件
var GetDatabaseSchemaCmd core.CommandObject = &getDatabaseSchemaCommandObject{}

/// 取得資料庫schema的指令物件
type getDatabaseSchemaCommandObject struct{}

func (x *getDatabaseSchemaCommandObject) PrintDoc() {
	args := newGetDatabaseSchemaArgs([]string{})
	args.PrintDoc()
}

func (x *getDatabaseSchemaCommandObject) Execute(subArgs []string) {
	args := newGetDatabaseSchemaArgs(subArgs)
	err := args.Parse(subArgs)
	if err != nil {
		args.PrintDoc()
	} else {
		runGetDatabaseSchemaData(args)
	}
}

/// new指令參數物件
func newGetDatabaseSchemaArgs(subArgs []string) *getDatabaseSchemaArgs {
	args := new(getDatabaseSchemaArgs)
	args.flagSet = flag.NewFlagSet("mssqlGetDatabaseSchema", flag.ExitOnError)
	args.GlobalConfig.Bind(args.flagSet)
	args.flagSet.StringVar(&args.ip, "ip", "", "(必填)ip")
	args.flagSet.StringVar(&args.port, "port", "", "(必填)port")
	args.flagSet.StringVar(&args.acc, "acc", "", "(必填)帳號")
	args.flagSet.StringVar(&args.pw, "pw", "", "(必填)密碼")
	args.flagSet.StringVar(&args.dbname, "dbname", "", "(必填)Database名稱")
	args.flagSet.BoolVar(&args.json, "json", true, "產生json格式，預設true")
	args.flagSet.BoolVar(&args.md, "md", false, "產生markdown格式，預設false")

	args.flagSet.Usage = args.PrintDoc
	return args
}

/// new指令參數物件
type getDatabaseSchemaArgs struct {
	flagSet *flag.FlagSet
	ip      string
	port    string
	acc     string
	pw      string
	dbname  string
	json    bool
	md      bool
	core.GlobalConfig
}

func (x *getDatabaseSchemaArgs) Parse(subArgs []string) (err error) {
	err = x.flagSet.Parse(subArgs)
	if err != nil {
		return
	}

	if x.ip == "" || x.port == "" || x.acc == "" || x.pw == "" || x.dbname == "" {
		err = errors.New("ip,port,acc,pw,dbname為必填")
	}
	return
}

func (x *getDatabaseSchemaArgs) PrintDoc() {
	doc := `
取得MSSQL的指定資料庫Schema資料

toolbox.exe mssqlGetDatabaseSchema -ip ip -port port -acc acc -pw pw -dbname dbname 
	`
	fmt.Println(doc)
	x.flagSet.PrintDefaults()
	fmt.Println("")
}

func runGetDatabaseSchemaData(args *getDatabaseSchemaArgs) {
	db, err := NewDb(args.ip, args.port, args.acc, args.pw, args.dbname)
	if err != nil {
		fmt.Println("connect to DB fail:", err.Error())
	}
	defer db.Close()

	databaseSchema, err := GetDatabaseSchema(db)
	if err != nil {
		fmt.Println("讀取DatabaseSchema失敗:", err.Error())
		return
	}

	if args.json {
		json, err := databaseSchema.ToJson()
		if err != nil {
			fmt.Println("json.Marshal失敗:", err.Error())
		}

		filepath := databaseSchema.Name + "資料庫Schema.json"
		err = tools.WriteTextFile(json, filepath)
		if err != nil {
			fmt.Println("寫入json失敗:", err.Error())
		}
	}

	if args.md {
		md := databaseSchema.ToMarkdown()

		filepath := databaseSchema.Name + "資料庫Schema.md"
		err := tools.WriteTextFile(md, filepath)
		if err != nil {
			fmt.Println("寫入md失敗:", err.Error())
		}
	}
}

func (x *DatabaseSchema) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString("\n# " + x.Name + "資料庫Schema \n")
	sb.WriteString("## Tables \n")
	for _, ts := range x.Tables {
		sb.WriteString(tableSchema2MdRow(ts))
	}

	sb.WriteString("## Views \n")
	for _, vs := range x.Views {
		sb.WriteString(viewSchema2MdRow(vs))
	}
	return sb.String()
}

func tableSchema2MdRow(x *TableSchema) string {
	var sb strings.Builder
	sb.WriteString("\n### " + x.Name + "\n")
	sb.WriteString(x.Description + "\n")
	sb.WriteString("| ColumnName | Nullable | Datatype | Description | PrimaryKey | \n")
	sb.WriteString("| --- | :-: | --- | --- | :-: | \n")
	for _, cs := range x.Columns {
		sb.WriteString(columnSchema2MdRow(cs) + "\n")
	}
	return sb.String()
}

func viewSchema2MdRow(x *ViewSchema) string {
	var sb strings.Builder
	sb.WriteString("\n### " + x.Name + "\n")
	sb.WriteString(x.Description + "\n")
	sb.WriteString("| ColumnName | Nullable | Datatype | Description | PrimaryKey | \n")
	sb.WriteString("| --- | :-: | --- | --- | :-: |\n")
	for _, cs := range x.Columns {
		sb.WriteString(columnSchema2MdRow(cs) + "\n")
	}
	return sb.String()
}

func columnSchema2MdRow(x *ColumnSchema) string {
	return fmt.Sprintf("|%s|%v|%s|%s|%v|", x.ColumnName, x.Nullable, x.Datatype, x.Description, x.IsPrimaryKey)
}
