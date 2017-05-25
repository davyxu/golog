package golog

import "os"

var logMap = make(map[string]*Logger)

func add(l *Logger) {

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

func VisitLogger(name string, callback func(*Logger) bool) {

	if name == "*" {

		for _, l := range logMap {
			if !callback(l) {
				break
			}
		}

	} else {
		l, ok := logMap[name]
		if !ok {
			return
		}

		if callback(l) {
			return
		}

	}
}

// 通过字符串设置某一类日志的级别
func SetLevelByString(loggerName string, level string) {

	VisitLogger(loggerName, func(l *Logger) bool {
		l.SetLevelByString(level)
		return true
	})

}

// 通过字符串设置某一类日志的崩溃级别
func SetPanicLevelByString(loggerName string, level string) {

	VisitLogger(loggerName, func(l *Logger) bool {
		l.SetPanicLevelByString(level)
		return true
	})
}

func SetColorFile(loggerName string, colorFileName string) {

	cf := NewColorFile()

	if err := cf.Load(colorFileName); err != nil {
		panic(err)
	}

	VisitLogger(loggerName, func(l *Logger) bool {
		l.SetColorFile(cf)
		return true
	})
}

func EnableColorLogger(loggerName string, enable bool) {

	VisitLogger(loggerName, func(l *Logger) bool {
		l.enableColor = enable
		return true
	})
}

func SetOutputLogger(loggerName string, filename string) {

	mode := os.O_RDWR | os.O_CREATE | os.O_APPEND

	f, err := os.OpenFile(filename, mode, 0666)
	if err != nil {
		panic(err)
	}

	VisitLogger(loggerName, func(l *Logger) bool {
		l.fileOutput = f
		return true
	})
}
