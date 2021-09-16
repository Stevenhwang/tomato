package utils

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
)

var red func(format string, a ...interface{})
var green func(format string, a ...interface{})
var yellow func(format string, a ...interface{})

func init() {
	red = color.New(color.FgRed).PrintfFunc()
	green = color.New(color.FgGreen).PrintfFunc()
	yellow = color.New(color.FgYellow).PrintfFunc()
}

func PrintRed(format string, a ...interface{}) {
	red(format, a...)
}

func PrintGreen(format string, a ...interface{}) {
	green(format, a...)
}

func PrintYellow(format string, a ...interface{}) {
	yellow(format, a...)
}

type ResultPrinter struct {
	SM sync.Mutex
}

func (rp *ResultPrinter) PrintFail(host string, format string, a ...interface{}) {
	rp.SM.Lock()
	color.Set(color.FgRed)
	fmt.Printf(
		"--------失败--------\n"+
			"主机: %s\n"+
			"结果: \n"+
			"\t"+fmt.Sprintf(format, a...)+"\n"+
			"-------------------\n", host,
	)
	color.Unset()
	rp.SM.Unlock()
}

func (rp *ResultPrinter) PrintSucc(host string, format string, a ...interface{}) {
	rp.SM.Lock()
	color.Set(color.FgGreen)
	fmt.Printf(
		"--------成功--------\n"+
			"主机: %s\n"+
			"结果: \n"+
			"\t"+fmt.Sprintf(format, a...)+"\n"+
			"-------------------\n", host,
	)
	color.Unset()
	rp.SM.Unlock()
}
