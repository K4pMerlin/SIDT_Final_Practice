package logger

import (
	"fmt"
	"strings"
)

func colorPrint(msg string, v ...any) {
	linuxColorPrintln(msg, fmt.Sprint(v...))
}

func colorPrintf(format string, msg string, v ...any) {
	// 上一行参数传v表示整体当成数组传入，参数传v... 表示多个参数分别传入
	linuxColorPrintln(msg, fmt.Sprintf(format, v...))
}

func linuxColorPrintln(msg string, str string) {
	color := logColorMap[strings.ToLower(msg)]
	myConsoleLogger.SetPrefix(color + "[" + msg + "]" + colorReset)
	overridePrintln(myConsoleLogger, logMap[logLevel] >= logMap[msg], color+str+colorReset)

	// 移除文件日志的颜色前缀
	myFileLogger.SetPrefix("[" + msg + "]")
	overridePrintln(myFileLogger, logMap[logLevel] >= logMap[msg], str)

}
