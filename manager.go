package golog

import (
	"io/ioutil"
	"os"
	"regexp"
	"sync"
)

var (
	loggerByName      = map[string]*Logger{}
	loggerByNameGuard sync.RWMutex
)

func add(l *Logger) {

	loggerByNameGuard.Lock()

	if _, ok := loggerByName[l.name]; ok {
		panic("duplicate logger name:" + l.name)
	}

	loggerByName[l.name] = l

	loggerByNameGuard.Unlock()
}

func str2loglevel(level string) Level {

	switch level {
	case "debug":
		return Level_Debug
	case "info":
		return Level_Info
	case "warn":
		return Level_Warn
	case "error", "err":
		return Level_Error
	case "fatal":
		return Level_Fatal
	}

	return Level_Debug
}

// 支持正则表达式查找logger， a|b|c指定多个日志, .表示所有日志
func VisitLogger(names string, callback func(*Logger) bool) error {

	loggerByNameGuard.RLock()

	defer loggerByNameGuard.RUnlock()

	exp, err := regexp.Compile(names)
	if err != nil {
		return err
	}

	for _, l := range loggerByName {

		if exp.MatchString(l.Name()) {
			if !callback(l) {
				break
			}
		}

	}

	return nil
}

// 通过字符串设置某一类日志的级别
func SetLevelByString(loggerName string, level string) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetLevelByString(level)
		return true
	})
}

// 通过字符串设置某一类日志的崩溃级别
func SetPanicLevelByString(loggerName string, level string) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetPanicLevelByString(level)
		return true
	})
}

func SetColorFile(loggerName string, colorFileName string) error {

	data, err := ioutil.ReadFile(colorFileName)
	if err != nil {
		return err
	}

	return SetColorDefine(loggerName, string(data))
}

func SetColorDefine(loggerName string, jsonFormat string) error {

	cf := NewColorFile()

	if err := cf.Load(jsonFormat); err != nil {
		return err
	}

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetColorFile(cf)
		return true
	})
}

func EnableColorLogger(loggerName string, enable bool) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.enableColor = enable
		return true
	})
}

func SetOutputLogger(loggerName string, filename string) error {

	mode := os.O_RDWR | os.O_CREATE | os.O_APPEND

	f, err := os.OpenFile(filename, mode, 0666)
	if err != nil {
		return err
	}

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.fileOutput = f
		return true
	})
}

func ClearAll() {

	loggerByNameGuard.Lock()
	loggerByName = map[string]*Logger{}
	loggerByNameGuard.Unlock()
}
