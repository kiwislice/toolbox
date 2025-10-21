package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/kiwislice/toolbox/mssql"
	"github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

var mssqlCompareTableDataCmd = &cobra.Command{
	Use:   "mssqlCompareTableData <setting file path>",
	Short: "Compare two MSSQL tables",
	Long:  `Compare two MSSQL tables based on a JSON configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		settingfile := args[0]
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
	},
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

func newDbL(args *settingArgs) (db *mssql.CustomDB, err error) {
	return mssql.NewDb(args.IpL, args.PortL, args.AccL, args.PwL, args.DbnameL)
}

func newDbR(args *settingArgs) (db *mssql.CustomDB, err error) {
	return mssql.NewDb(args.IpR, args.PortR, args.AccR, args.PwR, args.DbnameR)
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
func compareTableData(dbL, dbR *mssql.CustomDB, tablename string) (same bool) {
	listL, err := mssql.SelectAll(dbL, tablename)
	if err != nil {
		tools.Error("SelectAll fail: " + err.Error())
		return false
	}
	listR, err := mssql.SelectAll(dbR, tablename)
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

func init() {
	rootCmd.AddCommand(mssqlCompareTableDataCmd)
}
