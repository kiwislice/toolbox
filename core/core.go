package core

import (
	"flag"
)

/// 指令物件
type CommandObject interface {
	/// 顯示說明文件
	PrintDoc()
	/// 執行
	Execute(subArgs []string)
}

/// 通用參數
type GlobalConfig struct {
	Debug bool
}

func (x *GlobalConfig) Bind(flagSet *flag.FlagSet) {
	flagSet.BoolVar(&x.Debug, "debug", false, "顯示debug訊息")
}
