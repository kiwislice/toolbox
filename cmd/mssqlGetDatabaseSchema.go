package cmd

import (
	"fmt"
	"strings"

	"github.com/kiwislice/toolbox/mssql"
	"github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

var (
	ip     string
	port   string
	acc    string
	pw     string
	dbname string
	generateJson   bool
	md     bool
)

var mssqlGetDatabaseSchemaCmd = &cobra.Command{
	Use:   "mssqlGetDatabaseSchema",
	Short: "Get MSSQL database schema",
	Long:  `Get the schema of a specified MSSQL database.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := mssql.NewDb(ip, port, acc, pw, dbname)
		if err != nil {
			fmt.Println("connect to DB fail:", err.Error())
		}
		defer db.Close()

		databaseSchema, err := mssql.GetDatabaseSchema(db)
		if err != nil {
			fmt.Println("讀取DatabaseSchema失敗:", err.Error())
			return
		}

		if generateJson {
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

		if md {
			md := toMarkdown(databaseSchema)

			filepath := databaseSchema.Name + "資料庫Schema.md"
			err := tools.WriteTextFile(md, filepath)
			if err != nil {
				fmt.Println("寫入md失敗:", err.Error())
			}
		}
	},
}

func toMarkdown(schema *mssql.DatabaseSchema) string {
	var sb strings.Builder
	sb.WriteString("\n# " + schema.Name + "資料庫Schema \n")
	sb.WriteString("## Tables \n")
	for _, ts := range schema.Tables {
		sb.WriteString(tableSchema2MdRow(ts))
	}

	sb.WriteString("## Views \n")
	for _, vs := range schema.Views {
		sb.WriteString(viewSchema2MdRow(vs))
	}
	return sb.String()
}

func tableSchema2MdRow(x *mssql.TableSchema) string {
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

func viewSchema2MdRow(x *mssql.ViewSchema) string {
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

func columnSchema2MdRow(x *mssql.ColumnSchema) string {
	return fmt.Sprintf("|%s|%v|%s|%s|%v|", x.ColumnName, x.Nullable, x.Datatype, x.Description, x.IsPrimaryKey)
}

func init() {
	rootCmd.AddCommand(mssqlGetDatabaseSchemaCmd)
	mssqlGetDatabaseSchemaCmd.Flags().StringVar(&ip, "ip", "", "ip address")
	mssqlGetDatabaseSchemaCmd.Flags().StringVar(&port, "port", "", "port")
	mssqlGetDatabaseSchemaCmd.Flags().StringVar(&acc, "acc", "", "account")
	mssqlGetDatabaseSchemaCmd.Flags().StringVar(&pw, "pw", "", "password")
	mssqlGetDatabaseSchemaCmd.Flags().StringVar(&dbname, "dbname", "", "database name")
	mssqlGetDatabaseSchemaCmd.Flags().BoolVar(&generateJson, "json", true, "generate json format")
	mssqlGetDatabaseSchemaCmd.Flags().BoolVar(&md, "md", false, "generate markdown format")

	mssqlGetDatabaseSchemaCmd.MarkFlagRequired("ip")
	mssqlGetDatabaseSchemaCmd.MarkFlagRequired("port")
	mssqlGetDatabaseSchemaCmd.MarkFlagRequired("acc")
	mssqlGetDatabaseSchemaCmd.MarkFlagRequired("pw")
	mssqlGetDatabaseSchemaCmd.MarkFlagRequired("dbname")
}
