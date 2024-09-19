package logger

import (
	"fmt"
	"strings"
	"sync"
	"syscall"
)

func colorPrint(msg string, v ...any) {
	winColorPrintln(msg, fmt.Sprint(v...))

}
func colorPrintf(format string, msg string, v ...any) {
	winColorPrintln(msg, fmt.Sprintf(format, v...))
}

// windows 操作系统下调用系统api实现cmd着色
func winColorPrintln(msg string, str string) {
	var mu sync.Mutex

	// 加锁防止颜色显示出现冲突问题
	mu.Lock()
	defer mu.Unlock()
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	switch strings.ToLower(msg) {
	case "error":
		_, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(4))
	case "warning":
		_, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(6))
	case "info":
		_, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(2))
	case "debug":
		_, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(5))
	default:
	}

	myConsoleLogger.SetPrefix("[" + msg + "]")
	overridePrintln(myConsoleLogger, logMap[logLevel] >= logMap[msg], str)

	// 不需要移除前缀
	myFileLogger.SetPrefix("[" + msg + "]")
	overridePrintln(myFileLogger, logMap[logLevel] >= logMap[msg], str)

	_, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
}
