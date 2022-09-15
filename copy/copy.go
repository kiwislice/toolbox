package copy

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	tools "github.com/kiwislice/toolbox/tools"

	"github.com/fatih/color"
)

type CopyCmd struct {
	args *CopyArgs
}

func (x *CopyCmd) PrintDoc() {
	x.args.PrintDoc()
}

type CopyArgs struct {
	flagSet *flag.FlagSet
	src     string
	dest    string
	tools.GlobalConfig
}

func (x *CopyArgs) Parse(subArgs []string) (err error) {
	if len(subArgs) < 2 {
		return errors.New("參數長度必須大於等於2")
	}
	x.src = subArgs[0]
	x.dest = subArgs[1]
	return x.flagSet.Parse(subArgs[2:])
}

func (x *CopyArgs) PrintDoc() {
	doc := `
複製檔案或資料夾

toolbox.exe copy <source> <destination>
	`
	fmt.Println(doc)
	x.flagSet.PrintDefaults()
	fmt.Println("")
}

func newCopyArgs(subArgs []string) *CopyArgs {
	args := new(CopyArgs)
	args.flagSet = flag.NewFlagSet("copy", flag.ExitOnError)
	args.GlobalConfig.Bind(args.flagSet)
	// args.flagSet.StringVar(&args.src, "enable", "false", "enableaaa")

	args.flagSet.Usage = args.PrintDoc
	return args
}

func newCopyCmd(subArgs []string) (*CopyCmd, error) {
	args := newCopyArgs(subArgs)
	cmd := &CopyCmd{args}

	err := args.Parse(subArgs)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (cmd *CopyCmd) Run() {
	src := cmd.args.src
	dest := cmd.args.dest

	color.Cyan("開始複製：從 %s 到 %s", src, dest)
	defer color.Cyan("結束複製：從 %s 到 %s", src, dest)

	srcExist, srcInfo := tools.IsExist(src)
	if !srcExist {
		color.Red("來源不存在")
		return
	}
	destExist, destInfo := tools.IsExist(dest)

	srcType := tools.FileType(srcInfo)
	// 如果目標不存在則將目標當作跟來源同類
	destType := srcType
	if destExist {
		destType = tools.FileType(destInfo)
	}

	if srcType == tools.FILETYPE_UNKNOW || destType == tools.FILETYPE_UNKNOW {
		color.Red("來源或目的地類型不明")
		return
	}

	var function func(src, dest string) error
	if srcType == tools.FILETYPE_FILE {
		if destType == tools.FILETYPE_FILE {
			function = f2f
		} else {
			function = f2d
		}
	} else {
		if destType == tools.FILETYPE_FILE {
			function = d2f
		} else {
			function = d2d
		}
	}
	if err := function(src, dest); err != nil {
		color.Red(fmt.Sprint(err))
	}
}

func f2f(src, dest string) error {
	return tools.CopyFile(src, dest)
}

func f2d(src, dest string) error {
	filename := filepath.Base(src)
	dest = filepath.Join(dest, filename)
	return f2f(src, dest)
}

func d2f(src, dest string) error {
	return errors.New("資料夾無法複製到檔案")
}

func d2d(src, dest string) error {
	return tools.CopyDir(src, dest)
}

func Execute(subArgs []string) {
	cmd, err := newCopyCmd(subArgs)
	if err != nil {
		cmd.PrintDoc()
	} else {
		cmd.Run()
	}
}

func PrintDoc() {
	cmd, _ := newCopyCmd([]string{})
	cmd.PrintDoc()
}
