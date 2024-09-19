package logger

import (
	"CengkeHelper/setup"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorPurple = "\033[35m"

	colorReset = "\033[0m"
)

var myConsoleLogger *log.Logger
var myFileLogger *log.Logger
var logLevel string
var logMap = map[string]int{
	"error":   0,
	"warning": 1,
	"info":    2,
	"debug":   3,
}
var logColorMap = map[string]string{
	"error":   colorRed,
	"warning": colorYellow,
	"info":    colorGreen,
	"debug":   colorPurple,
}

// 仅显示时间，但写入文件
func init() {

	myConsoleLogger = log.New(os.Stdout, "[Default]", log.Lshortfile|log.Ltime)
	myFileLogger = log.New(os.Stdout, "[Default]", log.Lshortfile|log.Ltime)

	logLevel = setup.Config.LogLevel

}

func setFileLogger() *os.File {
	// fileLogger 的单独配置
	now := time.Now()
	//now.Format("2006-01-02")
	logDate := fmt.Sprintf("logs/%v_gin.log",
		now.Format("2006-01-02"))

	file, err := os.OpenFile(logDate, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	// 配置 fileLogger 的输出文件
	writer := io.MultiWriter(file)
	myFileLogger.SetOutput(writer)

	return file

}

// 重写 log 的Println 方法，修改调用堆栈的追踪深度，以便调试
func overridePrintln(l *log.Logger, isDisplay bool, str string) {
	if !isDisplay {
		return
	}

	file := setFileLogger()
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)

	// 修改了linux的调用追踪路径
	err := l.Output(5, str)
	if err != nil {
		return
	}
}

func Debug(v ...any) {
	colorPrint("Debug", v...)
}

func Info(v ...any) {
	colorPrint("Info", v...)
}

func Warning(v ...any) {
	colorPrint("Warning", v...)
}

func Error(v ...any) {
	colorPrint("Error", v...)
	//os.Exit(1)
}

// DebugF 带格式化的调试日志
func DebugF(format string, v ...any) {
	colorPrintf(format, "Debug", v...)
}

// InfoF 带格式化的信息日志
func InfoF(format string, v ...any) {
	colorPrintf(format, "Info", v...)
}

// WarningF 带格式化的警告日志
func WarningF(format string, v ...any) {
	colorPrintf(format, "Warning", v...)
}

// ErrorF 带格式化的错误日志
func ErrorF(format string, v ...any) {
	colorPrintf(format, "Error", v...)
	os.Exit(1)
}
