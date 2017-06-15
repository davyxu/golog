package golog

import (
	"errors"
	"os"
)

var logMap = make(map[string]*Logger)

func add(l *Logger) {

	if _, ok := logMap[l.name]; ok {
		panic("duplicate logger name:" + l.name)
	}

	logMap[l.name] = l
}

func str2loglevel(level string) Level {

	switch level {
	case "debug":
		return Level_Debug
	case "info":
		return Level_Info
	case "warn":
		return Level_Warn
	case "error":
		return Level_Error
	case "fatal":
		return Level_Fatal
	}

	return Level_Debug
}

var ErrLoggerNotFound = errors.New("logger not found")

func VisitLogger(name string, callback func(*Logger) bool) error {

	if name == "*" {

		for _, l := range logMap {
			if !callback(l) {
				break
			}
		}

	} else {
		l, ok := logMap[name]
		if !ok {
			return ErrLoggerNotFound
		}

		if callback(l) {
			return nil
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

	cf := NewColorFile()

	if err := cf.Load(colorFileName); err != nil {
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
