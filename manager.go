package golog

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

func selectLogger(name string, callback func(*Logger)) {

	if name == "*" {

		for _, l := range logMap {
			callback(l)
		}

	} else {
		l, ok := logMap[name]
		if !ok {
			return
		}

		callback(l)

	}
}

// 通过字符串设置某一类日志的级别
func SetLevelByString(loggerName string, level string) {

	selectLogger(loggerName, func(l *Logger) {
		l.SetLevelByString(level)
	})

}

// 通过字符串设置某一类日志的崩溃级别
func SetPanicLevelByString(loggerName string, level string) {

	selectLogger(loggerName, func(l *Logger) {
		l.SetPanicLevelByString(level)
	})
}

func SetColorFile(loggerName string, colorFileName string) {

	cf := NewColorFile()

	if err := cf.Load(colorFileName); err != nil {
		panic(err)
	}

	selectLogger(loggerName, func(l *Logger) {
		l.SetColorFile(cf)
	})
}
