package tools

import "github.com/fatih/color"

func Info(text string) {
	color.White(text)
}

func Infof(format string, args ...any) {
	color.White(format, args...)
}

func Warn(text string) {
	color.Yellow(text)
}

func Warnf(format string, args ...any) {
	color.Yellow(format, args...)
}

func Error(text string) {
	color.Red(text)
}

func Errorf(format string, args ...any) {
	color.Red(format, args...)
}

func Debug(text string) {
	color.Green(text)
}

func Debugf(format string, args ...any) {
	color.Green(format, args...)
}
