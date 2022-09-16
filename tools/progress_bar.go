package tools

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	anime1 = []string{
		" [=     ] ",
		" [ =    ] ",
		" [  =   ] ",
		" [   =  ] ",
		" [    = ] ",
		" [     =] ",
		" [    = ] ",
		" [   =  ] ",
		" [  =   ] ",
		" [ =    ] ",
	}
	anime2         = `|/-\`
	colorFunctions = []func(string, ...any) string{
		color.BlackString,
		color.BlueString,
		color.CyanString,
		color.GreenString,
		color.HiBlackString,
		color.HiBlueString,
		color.HiCyanString,
		color.HiGreenString,
		color.HiMagentaString,
		color.HiRedString,
		color.HiWhiteString,
		color.HiYellowString,
		color.MagentaString,
		color.RedString,
		color.WhiteString,
		color.YellowString,
	}
)

type LoadingText interface {
	SetText(text string)
	Start()
	Finish()
}

func NewLoadingText() LoadingText {
	return &animeLoadingText{animeTextProvider: anime1ColorfulTextProvider}
}

type simpleLoadingText struct {
	maxLen int
}

func (x *simpleLoadingText) SetText(text string) {
	x.maxLen = max(x.maxLen, len(text))
	if len(text) < x.maxLen {
		text += strings.Repeat(" ", x.maxLen-len(text))
	}
	fmt.Print("\r" + text)
}

func (x *simpleLoadingText) Start() {
}

func (x *simpleLoadingText) Finish() {
	fmt.Print("\r" + strings.Repeat(" ", x.maxLen) + "\r")
}

type animeTextProvider func(index int) string

func anime1TextProvider(index int) string {
	return anime1[index%len(anime1)]
}

func anime2TextProvider(index int) string {
	i := index % len(anime2)
	return " " + anime2[i:i+1] + " "
}

func anime1ColorfulTextProvider(index int) string {
	s := anime1[index%len(anime1)]
	return colorFunctions[index%len(colorFunctions)](s)
}

type animeLoadingText struct {
	text       string
	maxLen     int
	running    bool
	animeIndex int
	animeTextProvider
}

func (x *animeLoadingText) SetText(text string) {
	text = x.animeTextProvider(x.animeIndex) + text
	x.maxLen = max(x.maxLen, len(text))
	if len(text) < x.maxLen {
		text += strings.Repeat(" ", x.maxLen-len(text))
	}
	x.text = text
}

func (x *animeLoadingText) Start() {
	x.running = true
	go func() {
		for x.running {
			fmt.Print("\r" + x.text)
			x.animeIndex++
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (x *animeLoadingText) Finish() {
	x.running = false
	x.text = ""
	fmt.Print("\r" + strings.Repeat(" ", x.maxLen) + "\r")
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}
