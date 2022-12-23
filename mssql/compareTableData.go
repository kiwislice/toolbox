package mssql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
	core "github.com/kiwislice/toolbox/core"
	tools "github.com/kiwislice/toolbox/tools"
)

/// 比較資料表內容的指令物件
var CompareTableData core.CommandObject = &compareTableDataCommandObject{}

/// 比較資料表內容的指令物件
type compareTableDataCommandObject struct{}

func (x *compareTableDataCommandObject) PrintDoc() {
	args := newCompareTableDataArgs([]string{})
	args.PrintDoc()
}

func (x *compareTableDataCommandObject) Execute(subArgs []string) {
	args := newCompareTableDataArgs(subArgs)
	err := args.Parse(subArgs)
	if err != nil {
		args.PrintDoc()
	} else {
		run(args)
	}
}

/// new指令參數物件
func newCompareTableDataArgs(subArgs []string) *compareTableDataArgs {
	args := new(compareTableDataArgs)
	args.flagSet = flag.NewFlagSet("mssqlCompareTableData", flag.ExitOnError)
	args.GlobalConfig.Bind(args.flagSet)
	// args.flagSet.StringVar(&args.src, "enable", "false", "enableaaa")

	args.flagSet.Usage = args.PrintDoc
	return args
}

/// new指令參數物件
type compareTableDataArgs struct {
	flagSet     *flag.FlagSet
	settingfile string
	core.GlobalConfig
}

func (x *compareTableDataArgs) Parse(subArgs []string) (err error) {
	if len(subArgs) < 1 {
		return errors.New("參數長度必須大於等於1")
	} else {
		x.settingfile = subArgs[0]
		return x.flagSet.Parse(subArgs[1:])
	}
}

func (x *compareTableDataArgs) PrintDoc() {
	doc := `
比較MSSQL的2張資料表內容是否完全相等，只比較共同欄位的資料

toolbox.exe mssqlCompareTableData <設定檔路徑>

設定檔範例

{
    "ipL": "127.0.0.1",
    "portL": "1433",
    "accL": "sa",
    "pwL": "<YourStrong@Passw0rd>",
    "dbnameL": "CYBusNew1",
    "ipR": "127.0.0.1",
    "portR": "1433",
    "accR": "sa",
    "pwR": "<YourStrong@Passw0rd>",
    "dbnameR": "CYBusOldOld1",
    "tables": ["Cost_18","Cost_DriverSafe","ISP_Appeal"]
}
	`
	fmt.Println(doc)
	x.flagSet.PrintDefaults()
	fmt.Println("")
}

func run(args *compareTableDataArgs) {
	settingfile := args.settingfile
	color.Yellow(settingfile)

	settings, err := readSettingFile(settingfile)
	if err != nil {
		tools.Error("readSettingFile fail: " + err.Error())
	}

	dbL, err := newDbL(&settings)
	if err != nil {
		fmt.Println("connect to DB(Left) fail:", err.Error())
	}
	defer dbL.Close()

	dbR, err := newDbR(&settings)
	if err != nil {
		fmt.Println("connect to DB(Right) fail:", err.Error())
	}
	defer dbR.Close()

	color.Cyan("Compare Start：%s", settings)
	defer color.Cyan("Compare Finish：%s", settings)

	for _, table := range settings.Tables {
		startMsg := fmt.Sprintf("Table: %s Start", table)
		color.Green(startMsg)

		compareTableData(dbL, dbR, table)

		endMsg := fmt.Sprintf("Table: %s Finish", table)
		color.Green(endMsg)
	}
}

type settingArgs struct {
	IpL     string   `json:"ipL"`
	PortL   string   `json:"portL"`
	AccL    string   `json:"accL"`
	PwL     string   `json:"pwL"`
	DbnameL string   `json:"dbnameL"`
	IpR     string   `json:"ipR"`
	PortR   string   `json:"portR"`
	AccR    string   `json:"accR"`
	PwR     string   `json:"pwR"`
	DbnameR string   `json:"dbnameR"`
	Tables  []string `json:"tables"`
}

func newDbL(args *settingArgs) (db *sql.DB, err error) {
	return NewDb(args.IpL, args.PortL, args.AccL, args.PwL, args.DbnameL)
}

func newDbR(args *settingArgs) (db *sql.DB, err error) {
	return NewDb(args.IpR, args.PortR, args.AccR, args.PwR, args.DbnameR)
}

func readSettingFile(settingfile string) (setting settingArgs, err error) {
	content, err := ioutil.ReadFile(settingfile)
	if err != nil {
		tools.Error("ioutil.ReadFile fail: " + err.Error())
	}
	color.Yellow(string(content))
	err = json.Unmarshal(content, &setting)
	if err != nil {
		tools.Error("json.Unmarshal fail: " + err.Error())
	}
	return setting, err
}

/// true=都相同
func compareTableData(dbL, dbR *sql.DB, tablename string) (same bool) {
	listL, err := SelectAll(dbL, tablename)
	if err != nil {
		tools.Error("SelectAll fail: " + err.Error())
		return false
	}
	listR, err := SelectAll(dbR, tablename)
	if err != nil {
		tools.Error("SelectAll fail: " + err.Error())
		return false
	}

	same = true
	nrowL, nrowR := len(listL), len(listR)
	maxNrow := max(nrowL, nrowR)
	for i := 0; i < maxNrow; i++ {
		if i >= nrowL {
			tools.Errorf("Only DB(Right) has data: %s\n", listR[i])
			continue
		}
		if i >= nrowR {
			tools.Errorf("Only DB(Left) has data: %s\n", listL[i])
			continue
		}

		mapL, mapR := listL[i], listR[i]
		for key := range mapL {
			if valueR, exist := mapR[key]; exist {
				valueL := listL[i][key]
				if valueL != valueR {
					tools.Errorf("Left Right different: (Left)=%s，(Right)=%s\n", mapL, mapR)
					same = false
					break
				}
			}
		}
	}
	return same
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
