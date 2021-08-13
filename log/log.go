package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	// error为红色
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)

	// info为蓝色
	infoLog = log.New(os.Stdout, "\033[34m[error]\033[0m ", log.LstdFlags|log.Lshortfile)

	// Logger类型表示一个活动状态的记录日志的对象，它会生成一行行的输出写入一个io.Writer接口
	// 每一条日志操作会调用一次io.Writer接口的Write方法
	// Logger类型的对象可以被多个线程安全的同时使用，它会保证对io.Writer接口的顺序访问
	loggers = []*log.Logger{errorLog, infoLog}

	// 互斥锁
	mu sync.Mutex
)

// 被暴露的方法
var (
	Error  = errorLog.Println
	ErrorF = errorLog.Printf
	Info   = infoLog.Println
	InfoF  = infoLog.Printf
)

// 定义日志的层级
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// 设置日志的层级
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
