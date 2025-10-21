// 該檔案提供 log 相關的工具函式
package tools

import "github.com/fatih/color"

// Info 會印出白色的訊息
func Info(text string) {
	color.White(text)
}

// Infof 會印出白色的格式化訊息
func Infof(format string, args ...any) {
	color.White(format, args...)
}

// Warn 會印出黃色的訊息
func Warn(text string) {
	color.Yellow(text)
}

// Warnf 會印出黃色的格式化訊息
func Warnf(format string, args ...any) {
	color.Yellow(format, args...)
}

// Error 會印出紅色的訊息
func Error(text string) {
	color.Red(text)
}

// Errorf 會印出紅色的格式化訊息
func Errorf(format string, args ...any) {
	color.Red(format, args...)
}

// Debug 會印出綠色的訊息
func Debug(text string) {
	color.Green(text)
}

// Debugf 會印出綠色的格式化訊息
func Debugf(format string, args ...any) {
	color.Green(format, args...)
}
