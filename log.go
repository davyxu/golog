// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. It also has a predefined 'standard'
// Logger accessible through helper functions Print[f|ln], Fatal[f|ln], and
// Panic[f|ln], which are easier to use than creating a Logger manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.
package golog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu          sync.Mutex // ensures atomic writes; protects the following fields
	flag        LogFlag    // properties
	buf         []byte     // for accumulating text to write
	level       Level
	panicLevel  Level
	enableColor bool
	name        string
	colorFile   *ColorFile

	output io.Writer

	color Color
}

// New creates a new Logger.   The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.

func New(name string) *Logger {
	l := &Logger{
		flag:       LstdFlags,
		level:      Level_Debug,
		name:       name,
		panicLevel: Level_Fatal,
		output:     os.Stdout,
	}

	add(l)

	return l
}

func (self *Logger) SetFlag(v LogFlag) {
	self.flag = v
}

func (self *Logger) Flag() LogFlag {
	return self.flag
}

func (self *Logger) Name() string {
	return self.name
}

func (self *Logger) writeLevelAndName(level Level) {
	if self.flag.Contains(Llevel) {
		self.buf = append(self.buf, levelString[level]...)
		self.buf = append(self.buf, ' ')
	}

	if self.flag.Contains(Lname) {
		self.buf = append(self.buf, self.name...)
		self.buf = append(self.buf, ' ')
	}

}

func (self *Logger) selectColorByLevel(level Level) {

	if levelColor := colorFromLevel(level); levelColor != NoColor {
		self.color = levelColor
	}

}

func (self *Logger) selectColor(level Level, text string) {

	if self.enableColor && self.colorFile != nil && self.color == NoColor {
		self.color = self.colorFile.ColorFromText(text)
	}

	self.selectColorByLevel(level)

	return
}

func (self *Logger) Log(level Level, text string) {

	if level < self.level {
		return
	}

	self.buf = self.buf[:0]

	// 文本内容
	if self.color != NoColor {
		self.buf = append(self.buf, logColorPrefix[self.color]...)
	}

	self.writeLevelAndName(level)

	writeTimePart(self.flag, &self.buf)

	writeFilePart(self.flag, &self.buf)

	self.buf = append(self.buf, text...)

	// 颜色后缀
	if self.color != NoColor {
		self.buf = append(self.buf, logColorSuffix...)
	}

	// 回车
	if (len(text) > 0 && text[len(text)-1] != '\n') || len(text) == 0 {
		self.buf = append(self.buf, '\n')
	}

	self.output.Write(self.buf)

	if int(level) >= int(self.panicLevel) {
		panic(text)
	}

	self.color = NoColor

}

func (self *Logger) SetColor(name string) {
	self.color = colorByName[name]
}

func (self *Logger) Debugf(format string, v ...interface{}) {

	text := fmt.Sprintf(format, v...)

	self.selectColor(Level_Debug, text)

	self.Log(Level_Debug, text)
}

func (self *Logger) Debugln(v ...interface{}) {

	text := fmt.Sprint(v...)

	self.selectColor(Level_Debug, text)

	self.Log(Level_Debug, text)
}

func (self *Logger) Infof(format string, v ...interface{}) {

	text := fmt.Sprintf(format, v...)

	self.selectColor(Level_Info, text)

	self.Log(Level_Info, text)
}

func (self *Logger) Infoln(v ...interface{}) {

	text := fmt.Sprint(v...)

	self.selectColor(Level_Info, text)

	self.Log(Level_Info, text)
}

func (self *Logger) Warnf(format string, v ...interface{}) {

	text := fmt.Sprintf(format, v...)

	self.selectColor(Level_Warn, text)

	self.Log(Level_Warn, text)
}

func (self *Logger) Warnln(v ...interface{}) {

	text := fmt.Sprint(v...)

	self.selectColor(Level_Warn, text)

	self.Log(Level_Warn, text)
}

func (self *Logger) Errorf(format string, v ...interface{}) {

	text := fmt.Sprintf(format, v...)

	self.selectColor(Level_Error, text)

	self.Log(Level_Error, text)
}

func (self *Logger) Errorln(v ...interface{}) {

	text := fmt.Sprint(v...)

	self.selectColor(Level_Error, text)

	self.Log(Level_Error, text)
}

func (self *Logger) Fatalf(format string, v ...interface{}) {

	text := fmt.Sprintf(format, v...)

	self.selectColor(Level_Fatal, text)

	self.Log(Level_Fatal, text)
}

func (self *Logger) Fatalln(v ...interface{}) {

	text := fmt.Sprint(v...)

	self.selectColor(Level_Fatal, text)

	self.Log(Level_Fatal, text)
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

func (self *Logger) SetPanicLevelByString(level string) {
	self.panicLevel = str2loglevel(level)

}

// 注意, 加色只能在Gogland的main方式启用, Test方式无法加色
func (self *Logger) SetColorFile(file *ColorFile) {
	self.colorFile = file
}
func (self *Logger) IsDebugEnabled() bool {
	return self.level == Level_Debug
}
