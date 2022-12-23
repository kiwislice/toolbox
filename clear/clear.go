package clear

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/kiwislice/toolbox/core"
	tools "github.com/kiwislice/toolbox/tools"

	"github.com/fatih/color"
)

type ClearCmd struct {
	args *ClearArgs
}

func (x *ClearCmd) PrintDoc() {
	x.args.PrintDoc()
}

type ClearArgs struct {
	flagSet *flag.FlagSet
	target  string
	core.GlobalConfig
}

func (x *ClearArgs) Parse(subArgs []string) (err error) {
	if len(subArgs) < 1 {
		return errors.New("參數長度必須大於等於1")
	}
	x.target = subArgs[0]
	return x.flagSet.Parse(subArgs[1:])
}

func (x *ClearArgs) PrintDoc() {
	doc := `
清空資料夾

toolbox.exe clear <target>
	`
	fmt.Println(doc)
	x.flagSet.PrintDefaults()
	fmt.Println("")
}

func newClearArgs(subArgs []string) *ClearArgs {
	args := new(ClearArgs)
	args.flagSet = flag.NewFlagSet("clear", flag.ExitOnError)
	args.GlobalConfig.Bind(args.flagSet)
	// args.flagSet.StringVar(&args.src, "enable", "false", "enableaaa")

	args.flagSet.Usage = args.PrintDoc
	return args
}

func newClearCmd(subArgs []string) (*ClearCmd, error) {
	args := newClearArgs(subArgs)
	cmd := &ClearCmd{args}

	err := args.Parse(subArgs)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (cmd *ClearCmd) Run() {
	target := cmd.args.target

	color.Cyan("開始清空：%s", target)
	defer color.Cyan("結束清空：%s", target)

	exist, fileInfo := tools.IsExist(target)
	if !exist || !fileInfo.IsDir() {
		tools.Warnf("目標不是資料夾：%s", target)
		return
	}

	err := os.RemoveAll(target)
	if err != nil {
		tools.Errorf("%s", err)
	}
	err = os.Mkdir(target, os.ModePerm)
	if err != nil {
		tools.Errorf("%s", err)
	}
}

func Execute(subArgs []string) {
	cmd, err := newClearCmd(subArgs)
	if err != nil {
		cmd.PrintDoc()
	} else {
		cmd.Run()
	}
}

func PrintDoc() {
	cmd, _ := newClearCmd([]string{})
	cmd.PrintDoc()
}
