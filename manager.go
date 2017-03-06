package golog

var logMap = make(map[string]*Logger)

func add(l *Logger) {

	logMap[l.name] = l
}

func str2loglevel(level string) int {

	switch level {
	case "debug":
		return LEVEL_DEBUG
	case "info":
		return LEVEL_INFO
	case "warn":
		return LEVEL_WARN
	case "error":
		return LEVEL_ERROR
	case "fatal":
		return LEVEL_FATAL
	}

	return LEVEL_DEBUG
}

func selectLogger(name string, callback func(*Logger)) {

	if name == "all" {

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
func SetLevelByString(name string, level string) {

	selectLogger(name, func(l *Logger) {
		l.SetLevelByString(level)
	})

}

// 通过字符串设置某一类日志的崩溃级别
func SetPanicLevelByString(name string, level string) {

	selectLogger(name, func(l *Logger) {
		l.SetPanicLevelByString(level)
	})
}
