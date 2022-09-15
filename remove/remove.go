package remove

import (
	"errors"
	"flag"
	"fmt"
	"os"

	tools "github.com/kiwislice/toolbox/tools"

	"github.com/fatih/color"
)

type RemoveCmd struct {
	args *RemoveArgs
}

func (x *RemoveCmd) PrintDoc() {
	x.args.PrintDoc()
}

type RemoveArgs struct {
	flagSet *flag.FlagSet
	target  string
	tools.GlobalConfig
}

func (x *RemoveArgs) Parse(subArgs []string) (err error) {
	if len(subArgs) < 1 {
		return errors.New("參數長度必須大於等於1")
	}
	x.target = subArgs[0]
	return x.flagSet.Parse(subArgs[1:])
}

func (x *RemoveArgs) PrintDoc() {
	doc := `
刪除檔案或資料夾

toolbox.exe remove <target>
	`
	fmt.Println(doc)
	x.flagSet.PrintDefaults()
	fmt.Println("")
}

func newRemoveArgs(subArgs []string) *RemoveArgs {
	args := new(RemoveArgs)
	args.flagSet = flag.NewFlagSet("remove", flag.ExitOnError)
	args.GlobalConfig.Bind(args.flagSet)
	// args.flagSet.StringVar(&args.src, "enable", "false", "enableaaa")

	args.flagSet.Usage = args.PrintDoc
	return args
}

func newRemoveCmd(subArgs []string) (*RemoveCmd, error) {
	args := newRemoveArgs(subArgs)
	cmd := &RemoveCmd{args}

	err := args.Parse(subArgs)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (cmd *RemoveCmd) Run() {
	target := cmd.args.target

	color.Cyan("開始刪除：%s", target)
	defer color.Cyan("結束刪除：%s", target)

	err := os.RemoveAll(target)
	if err != nil {
		tools.Errorf("%s", err)
	}
}

func Execute(subArgs []string) {
	cmd, err := newRemoveCmd(subArgs)
	if err != nil {
		cmd.PrintDoc()
	} else {
		cmd.Run()
	}
}

func PrintDoc() {
	cmd, _ := newRemoveCmd([]string{})
	cmd.PrintDoc()
}
