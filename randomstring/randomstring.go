package randomstring

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/kiwislice/toolbox/core"
	"github.com/kiwislice/toolbox/tools"
)

type RandomStringCmd struct {
	args *RandomStringArgs
}

func (x *RandomStringCmd) PrintDoc() {
	x.args.PrintDoc()
}

type RandomStringArgs struct {
	flagSet *flag.FlagSet
	src     string
	length  int
	count   int
	core.GlobalConfig
}

func (x *RandomStringArgs) Parse(subArgs []string) (err error) {
	return x.flagSet.Parse(subArgs)
}

func (x *RandomStringArgs) PrintDoc() {
	doc := `
產生隨機字串

toolbox.exe randomString
	`
	fmt.Println(doc)
	x.flagSet.PrintDefaults()
	fmt.Println("")
}

func newRandomStringArgs(subArgs []string) *RandomStringArgs {
	args := new(RandomStringArgs)
	args.flagSet = flag.NewFlagSet("randomString", flag.ExitOnError)
	args.GlobalConfig.Bind(args.flagSet)
	args.flagSet.StringVar(&args.src, "src", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", "字元來源")
	args.flagSet.IntVar(&args.length, "length", 8, "隨機字串長度")
	args.flagSet.IntVar(&args.count, "count", 1, "隨機字串數量")

	args.flagSet.Usage = args.PrintDoc
	return args
}

func newRandomStringCmd(subArgs []string) (*RandomStringCmd, error) {
	args := newRandomStringArgs(subArgs)
	cmd := &RandomStringCmd{args}

	err := args.Parse(subArgs)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (cmd *RandomStringCmd) Run() {
	tools.Debug("開始產生隨機字串")
	defer tools.Debug("結束產生隨機字串")

	for i := 0; i < cmd.args.count; i++ {
		s := RandomString(cmd.args.src, cmd.args.length)
		fmt.Println(s)
	}
}

func RandomString(src string, length int) string {
	seed := time.Now().UnixNano() + rand.Int63()
	r := rand.New(rand.NewSource(seed))
	ar := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		index := r.Intn(len(src))
		ar = append(ar, src[index])
	}
	return string(ar)
}

func Execute(subArgs []string) {
	cmd, err := newRandomStringCmd(subArgs)
	if err != nil {
		cmd.PrintDoc()
	} else {
		cmd.Run()
	}
}

func PrintDoc() {
	cmd, _ := newRandomStringCmd([]string{})
	cmd.PrintDoc()
}
