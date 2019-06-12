package golog

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type OutputFileOption struct {
	MaxFileSize int // 单个文件最大尺寸
}

type CreateLogFileEventFunc func(logFileName string)

var (
	globalWriter         io.Writer        = os.Stdout // 默认输出到标准输出
	fileopt              OutputFileOption             // 日志文件参数
	globalOutputFileName string                       // 输入的日志文件名
	logFileHandle        *os.File                     // 当前写入的日志文件句柄
	currFileIndex        int                          // 当前日志索引号
	onCreateLogFile      CreateLogFileEventFunc
)

// 获取当前日志索引号
func GetOutputLogIndex() int {
	return currFileIndex
}

// 获取设置的日志文件名
func GetOutputFileName() string {
	return globalOutputFileName
}

// 获取输出日志文件的参数
func GetOutputFileOption() OutputFileOption {
	return fileopt
}

// 设置全局日志输出到文件,
// 默认单文件模式: 可以忽略optList MaxFileSize=0
// 多文件模式: 每当日志超过MaxFileSize, 就会重新生成新的文件, 不修改已经生成的日志, 不是Rolling模式
func SetOutputToFile(filename string, optList ...interface{}) error {

	globalOutputFileName = filename

	// 自动创建日志目录
	logDir := filepath.Dir(globalOutputFileName)

	_, err := os.Stat(logDir)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(logDir, 0777)
	}

	for _, opt := range optList {
		switch v := opt.(type) {
		case OutputFileOption:
			fileopt = v

			// 如果有约束最大文件尺寸, 说明需要创建多个文件
			if v.MaxFileSize > 0 {

				// 如果max文件不存在, 创建max文件
				if !readMaxFile() {
					writeMaxFile()
				}
			}
		}
	}

	return nil
}

func writeMaxFile() {
	// 更新最大日志文件
	maxFile := fmt.Sprintf("%s.max", globalOutputFileName)

	maxLogIndex := strconv.Itoa(currFileIndex)

	ioutil.WriteFile(maxFile, []byte(maxLogIndex), 0666)
}

func readMaxFile() bool {

	maxFile := fmt.Sprintf("%s.max", globalOutputFileName)
	data, err := ioutil.ReadFile(maxFile)
	if err != nil {
		return false
	}

	v, err := strconv.Atoi(string(data))
	if err != nil {
		return false
	}

	currFileIndex = v

	return true
}

// 输出到自定义接口
func SetOutput(loggerName string, output io.Writer) error {
	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetOutptut(output)
		return true
	})
}

// 设置创建日志文件的回调
func SetCreateLogEvent(callback CreateLogFileEventFunc) {
	onCreateLogFile = callback
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

	// 不用创建文件
	if globalOutputFileName == "" {
		return
	}

	// 文件已经创建
	if logFileHandle != nil {

		// 设定了文件大小约束
		if fileopt.MaxFileSize > 0 {
			info, err := logFileHandle.Stat()
			if err != nil {
				return
			}

			// 文件还没到达约束
			if info.Size() < int64(fileopt.MaxFileSize) {
				return
			}
		} else { // 没有设置约束, 直接返回
			return
		}

	}

	// 找到可用的文件名
	fileName := foundUseableName(globalOutputFileName)

	// 关闭之前的日志文件
	if logFileHandle != nil {
		logFileHandle.Close()
	}

	// 首次创建/尺寸超过后新创建
	var err error
	logFileHandle, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		globalWriter = logFileHandle
	}

	// 创建新文件时, 更新max文件
	writeMaxFile()

	if onCreateLogFile != nil {
		onCreateLogFile(fileName)
	}

}

func foundUseableName(originName string) (ret string) {

	for i := currFileIndex; ; i++ {

		if i == 0 {
			ret = originName
		} else {
			ret = fmt.Sprintf("%s.%d", originName, i)
		}

		if _, err := os.Stat(ret); err != nil {
			if os.IsNotExist(err) {

				currFileIndex = i

				return
			}
		}

	}

}
