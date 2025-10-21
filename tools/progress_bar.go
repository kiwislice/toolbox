// 該檔案提供 terminal 進度條相關的工具函式
package tools

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	// anime1 是一個動畫的 frame a
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
	// anime2 是一個動畫的 frame b
	anime2 = `|/-\`
	// colorFunctions 是一個包含所有 color 函式的 slice
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

// LoadingText 是一個介面，定義了進度條需要有的方法
type LoadingText interface {
	SetText(text string)
	Start()
	Finish()
}

// NewLoadingText 會回傳一個新的 LoadingText
func NewLoadingText() LoadingText {
	return &animeLoadingText{animeTextProvider: anime1ColorfulTextProvider}
}

// simpleLoadingText 是一個簡單的進度條實作，只會顯示文字
type simpleLoadingText struct {
	maxLen int
}

// SetText 會設定進度條顯示的文字
func (x *simpleLoadingText) SetText(text string) {
	x.maxLen = max(x.maxLen, len(text))
	if len(text) < x.maxLen {
		text += strings.Repeat(" ", x.maxLen-len(text))
	}
	fmt.Print("\r" + text)
}

// Start 會開始進度條
func (x *simpleLoadingText) Start() {
}

// Finish 會結束進度條
func (x *simpleLoadingText) Finish() {
	fmt.Print("\r" + strings.Repeat(" ", x.maxLen) + "\r")
}

// animeTextProvider 是一個函式類型，用來提供動畫的 frame
type animeTextProvider func(index int) string

// anime1TextProvider 會回傳 anime1 的 frame
func anime1TextProvider(index int) string {
	return anime1[index%len(anime1)]
}

// anime2TextProvider 會回傳 anime2 的 frame
func anime2TextProvider(index int) string {
	i := index % len(anime2)
	return " " + anime2[i:i+1] + " "
}

// anime1ColorfulTextProvider 會回傳帶有顏色的 anime1 frame
func anime1ColorfulTextProvider(index int) string {
	s := anime1[index%len(anime1)]
	return colorFunctions[index%len(colorFunctions)](s)
}

// animeLoadingText 是一個帶有動畫的進度條實作
type animeLoadingText struct {
	text              string
	maxLen            int
	running           bool
	animeIndex        int
	animeTextProvider // 匿名字段，可以直接呼叫 animeTextProvider 的方法
}

// SetText 會設定進度條顯示的文字
func (x *animeLoadingText) SetText(text string) {
	text = x.animeTextProvider(x.animeIndex) + text
	x.maxLen = max(x.maxLen, len(text))
	if len(text) < x.maxLen {
		text += strings.Repeat(" ", x.maxLen-len(text))
	}
	x.text = text
}

// Start 會開始進度條
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

// Finish 會結束進度條
func (x *animeLoadingText) Finish() {
	x.running = false
	x.text = ""
	fmt.Print("\r" + strings.Repeat(" ", x.maxLen) + "\r")
}

// max 會回傳兩個整數中較大的那個
func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}
