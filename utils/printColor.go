package utils

import (
	"github.com/fatih/color"
)

func PrintRed(format string, a ...interface{}) {
	red := color.New(color.FgRed).PrintfFunc()
	red(format, a...)
}

func PrintGreen(format string, a ...interface{}) {
	green := color.New(color.FgGreen).PrintfFunc()
	green(format, a...)
}
