/*
	对log包的二次封装
	1、可进行分级；
	2、选择输出目的地；

	bug:
	1、日志只能输出到一个目的地；
*/
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	LEVEL_FATAL = iota
	LEVEL_ERROR
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_DEBUG
)

type Logger struct {
	level  int         //日志等级
	logger *log.Logger //
}

/*
	创建一个*Logger对象

	参数：
	level：输出日志等级
	dst：日志输出目的地，主要实现了io.Writer接口均可

	返回：
	*Logger
*/
func NewLogger(level int, dst io.Writer) *Logger {
	logger := log.New(dst, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		level:  level,
		logger: logger,
	}
}

/*
	输出日志到文件，若文件不存在会自动创建

	参数：
	level：日志等级
	path：日志文件路径

	返回：
	*Logger
*/
func NewLoggerWithFile(level int, path string) *Logger {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln("初始化日志失败：", err)
	}

	return NewLogger(level, file)
}

func (l *Logger) DEBUG(arg ...interface{}) {
	if l.level >= LEVEL_DEBUG {
		l.logger.SetPrefix(fmt.Sprintf("%-6s", "DEBUG"))
		l.logger.Println(arg)
	}
}

func (l *Logger) INFO(arg ...interface{}) {
	if l.level >= LEVEL_INFO {
		l.logger.SetPrefix(fmt.Sprintf("%-6s", "INFO"))
		l.logger.Println(arg)
	}
}

func (l *Logger) WARN(arg ...interface{}) {
	if l.level >= LEVEL_WARN {
		l.logger.SetPrefix(fmt.Sprintf("%-6s", "WARN"))
		l.logger.Println(arg)
	}
}

func (l *Logger) ERROR(arg ...interface{}) {
	if l.level >= LEVEL_ERROR {
		l.logger.SetPrefix(fmt.Sprintf("%-6s", "ERROR"))
		l.logger.Println(arg)
	}
}

func (l *Logger) FATAL(arg ...interface{}) {
	if l.level >= LEVEL_FATAL {
		l.logger.SetPrefix(fmt.Sprintf("%-6s", "FATAL"))
		l.logger.Fatalln(arg)
	}
}
