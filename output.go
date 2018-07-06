package golog

import (
	"io"
	"os"
)

// 输出到文件
func SetOutputToFile(loggerName string, filename string) error {

	mode := os.O_RDWR | os.O_CREATE | os.O_APPEND

	fileHandle, err := os.OpenFile(filename, mode, 0666)
	if err != nil {
		return err
	}

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetOutptut(fileHandle)
		return true
	})
}

// 输出到自定义接口
func SetOutput(loggerName string, output io.Writer) error {
	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetOutptut(output)
		return true
	})
}
