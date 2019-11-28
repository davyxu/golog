package golog

import "fmt"

type Level int

const (
	Level_Debug Level = iota
	Level_Info
	Level_Warn
	Level_Error
)

var levelString = [...]string{
	"[DEBU]",
	"[INFO]",
	"[WARN]",
	"[ERRO]",
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

func (self *Logger) Debugf(format string, v ...interface{}) {

	self.LogText(Level_Debug, fmt.Sprintf(format, v...), nil)
}

func (self *Logger) Infof(format string, v ...interface{}) {

	self.LogText(Level_Info, fmt.Sprintf(format, v...), nil)
}

func (self *Logger) Warnf(format string, v ...interface{}) {

	self.LogText(Level_Warn, fmt.Sprintf(format, v...), nil)
}

func (self *Logger) Errorf(format string, v ...interface{}) {

	self.LogText(Level_Error, fmt.Sprintf(format, v...), nil)
}

func (self *Logger) Debugln(v ...interface{}) {

	self.LogText(Level_Debug, fmt.Sprintln(v...), nil)
}

func (self *Logger) Infoln(v ...interface{}) {

	self.LogText(Level_Info, fmt.Sprintln(v...), nil)
}

func (self *Logger) Warnln(v ...interface{}) {

	self.LogText(Level_Warn, fmt.Sprintln(v...), nil)
}

func (self *Logger) Errorln(v ...interface{}) {
	self.LogText(Level_Error, fmt.Sprintln(v...), nil)
}

func (self *Logger) SetLevelByString(level string) {

	self.SetLevel(str2loglevel(level))

}

func (self *Logger) SetLevel(lv Level) {
	self.level = lv
}

func (self *Logger) Level() Level {
	return self.level
}

func (self *Logger) CurrLevelString() string {
	return levelString[self.currLevel]
}
