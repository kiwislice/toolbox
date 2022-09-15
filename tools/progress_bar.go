package tools

import (
	"fmt"
	"strings"
	"time"
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
	anime2 = `|/-\`
)

type LoadingText interface {
	SetText(text string)
	Start()
	Finish()
}

func NewLoadingText() LoadingText {
	return &animeLoadingText{animeTextProvider: anime2TextProvider}
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
			time.Sleep(20 * time.Millisecond)
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
