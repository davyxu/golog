package golog

import (
	"io/ioutil"
	"strings"
)

type Color int

const (
	NoColor Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Purple
	DarkGreen
	White
)

var logColorPrefix = []string{
	"",
	"\x1b[030m",
	"\x1b[031m",
	"\x1b[032m",
	"\x1b[033m",
	"\x1b[034m",
	"\x1b[035m",
	"\x1b[036m",
	"\x1b[037m",
}

type colorData struct {
	name string
	c    Color
}

var colorByName = []colorData{
	{"none", NoColor},
	{"black", Black},
	{"red", Red},
	{"green", Green},
	{"yellow", Yellow},
	{"blue", Blue},
	{"purple", Purple},
	{"darkgreen", DarkGreen},
	{"white", White},
}

func matchColor(name string) Color {

	lower := strings.ToLower(name)

	for _, d := range colorByName {

		if d.name == lower {
			return d.c
		}
	}

	return NoColor
}

func colorFromLevel(l Level) Color {
	switch l {
	case Level_Warn:
		return Yellow
	case Level_Error:
		return Red
	}

	return NoColor
}

var logColorSuffix = "\x1b[0m"

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
		l.EnableColor(enable)
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

func (self *Logger) SetColor(name string) *Logger {
	self.mu.Lock()
	self.currColor = matchColor(name)
	self.mu.Unlock()

	return self
}

// 注意, 加色只能在Gogland的main方式启用, Test方式无法加色
func (self *Logger) SetColorFile(file *ColorFile) {
	self.colorFile = file
}

func (self *Logger) selectColorByLevel() {

	if levelColor := colorFromLevel(self.currLevel); levelColor != NoColor {
		self.currColor = levelColor
	}

}

func (self *Logger) selectColorByText() {

	if self.enableColor && self.colorFile != nil && self.currColor == NoColor {
		self.currColor = self.colorFile.ColorFromText(self.currText)
	}

	return
}
