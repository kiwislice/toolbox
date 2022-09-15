package tools

import "flag"

type GlobalConfig struct {
	Debug bool
}

func (x *GlobalConfig) Bind(flagSet *flag.FlagSet) {
	flagSet.BoolVar(&x.Debug, "debug", false, "顯示debug訊息")
}
