package golog

type Level int

const (
	Level_Debug Level = iota
	Level_Info
	Level_Warn
	Level_Error
	Level_Fatal
)

var levelString = [...]string{
	"[DEBU]",
	"[INFO]",
	"[WARN]",
	"[ERRO]",
	"[FATL]",
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
