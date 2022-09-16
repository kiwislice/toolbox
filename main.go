package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kiwislice/toolbox/clear"
	"github.com/kiwislice/toolbox/copy"
	"github.com/kiwislice/toolbox/randomstring"
	"github.com/kiwislice/toolbox/remove"
	"github.com/kiwislice/toolbox/tools"
)

func printMainDoc() {
	doc := `
自製工具箱

toolbox.exe <command> [<args>]

以下為<command>清單：
	help	顯示使用說明
	copy	複製檔案或資料夾
	remove	刪除檔案或資料夾

'toolbox.exe help <command>'可以看到該command的使用說明
	`
	fmt.Println(doc)
}

func main() {
	if len(os.Args) <= 1 {
		printMainDoc()
		return
	}

	needHelp := os.Args[1] == "help"
	command := os.Args[1]
	args := os.Args[2:]
	if needHelp && len(os.Args) >= 3 {
		command = os.Args[2]
		args = os.Args[3:]
	}

	switch command {
	case "copy":
		if needHelp {
			copy.PrintDoc()
		} else {
			copy.Execute(args)
		}
	case "remove":
		if needHelp {
			remove.PrintDoc()
		} else {
			remove.Execute(args)
		}
	case "clear":
		if needHelp {
			clear.PrintDoc()
		} else {
			clear.Execute(args)
		}
	case "randomString":
		if needHelp {
			randomstring.PrintDoc()
		} else {
			randomstring.Execute(args)
		}
	case "testLoadingText":
		pb := tools.NewLoadingText()
		pb.Start()
		for i := 1 << 30; i > 1; i = int(float64(i) / 1.1) {
			s := fmt.Sprintf("aaaaaaaaa %d bbbbb", i)
			pb.SetText(s)
			time.Sleep(10 * time.Millisecond)
		}
		pb.Finish()
	case "showText":
		fmt.Print(args)

	default:
		printMainDoc()
	}
}
