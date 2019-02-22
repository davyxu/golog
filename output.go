package golog

import (
	"fmt"
	"io"
	"os"
)

type OutputFileOption struct {
	MaxFileSize int // 单个文件最大尺寸
}

var (
	globalWriter         io.Writer = os.Stdout // 默认输出到标准输出
	fileopt              OutputFileOption
	globalOutputFileName string
)

// 设置全局日志输出到文件,
// 默认单文件模式: 可以忽略optList MaxFileSize=0
// 多文件模式: 每当日志超过MaxFileSize, 就会重新生成新的文件, 不修改已经生成的日志, 不是Rolling模式
func SetOutputToFile(filename string, optList ...interface{}) error {

	globalOutputFileName = filename

	for _, opt := range optList {
		switch v := opt.(type) {
		case OutputFileOption:
			fileopt = v
		}
	}

	return nil
}

// 输出到自定义接口
func SetOutput(loggerName string, output io.Writer) error {
	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetOutptut(output)
		return true
	})
}

// 设置本日志的输出, 单个日志设置输出时, 优先度高于全局
func (self *Logger) SetOutptut(writer io.Writer) {
	self.output = writer
}

func (self *Logger) GetOutput() io.Writer {
	return self.output
}

func globalWrite(b []byte) {

	checkFileCreate()

	if asyncWriteBuff != nil {
		queuedWrite(b)
	} else {

		globalWriter.Write(b)
	}
}

func checkFileCreate() {

	if globalWriter == os.Stdout {

		if globalOutputFileName == "" {
			return
		}

	} else {

		if fileopt.MaxFileSize > 0 {

			f := globalWriter.(*os.File)
			info, err := f.Stat()
			if err != nil {
				return
			}

			if info.Size() < int64(fileopt.MaxFileSize) {
				return
			}
		} else {
			return
		}
	}

	fileName := foundUseableName(globalOutputFileName)

	fileHandle, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		globalWriter = fileHandle
	}

}

func foundUseableName(name string) string {

	originName := name
	for i := 1; ; i++ {
		if _, err := os.Stat(name); err != nil {
			if os.IsNotExist(err) {
				return name
			}
		}

		name = fmt.Sprintf("%s.%d", originName, i)
	}

}
