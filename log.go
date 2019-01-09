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

type PartFunc func(*Logger)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu          sync.Mutex // ensures atomic writes; protects the following fields
	buf         []byte     // for accumulating text to write
	level       Level
	enableColor bool
	name        string
	colorFile   *ColorFile

	parts []PartFunc

	output io.Writer

	currColor     Color
	currLevel     Level
	currText      string
	currCondition bool
}

// New creates a new Logger.   The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.

const lineBuffer = 32

func New(name string) *Logger {
	l := &Logger{
		level:         Level_Debug,
		name:          name,
		output:        os.Stdout,
		buf:           make([]byte, 0, lineBuffer),
		currCondition: true,
	}

	l.SetParts(LogPart_Level, LogPart_Name, LogPart_Time)

	add(l)

	return l
}

func (self *Logger) SetOutptut(writer io.Writer) {
	self.output = writer
}

func (self *Logger) GetOutput() io.Writer {
	return self.output
}

func (self *Logger) EnableColor(v bool) {
	self.mu.Lock()
	self.enableColor = v
	self.mu.Unlock()
}

func (self *Logger) SetParts(f ...PartFunc) {

	self.parts = []PartFunc{logPart_ColorBegin}
	self.parts = append(self.parts, f...)
	self.parts = append(self.parts, logPart_Text, logPart_ColorEnd, logPart_Line)
}

// 二次开发接口
func (self *Logger) WriteRawString(s string) {
	self.buf = append(self.buf, s...)
}

func (self *Logger) WriteRawByte(b byte) {
	self.buf = append(self.buf, b)
}

func (self *Logger) Name() string {
	return self.name
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

func (self *Logger) Buff() []byte {
	return self.buf
}

func (self *Logger) Log(level Level, text string) {

	// 防止日志并发打印导致的文本错位
	self.mu.Lock()
	defer self.mu.Unlock()

	self.currLevel = level
	self.currText = text

	defer self.resetState()

	if self.currLevel < self.level || !self.currCondition {
		return
	}

	self.selectColorByText()
	self.selectColorByLevel()

	self.buf = self.buf[:0]

	for _, p := range self.parts {
		p(self)
	}

	self.output.Write(self.buf)
}

func (self *Logger) Condition(value bool) *Logger {

	self.mu.Lock()
	self.currCondition = value
	self.mu.Unlock()

	return self
}

func (self *Logger) resetState() {
	self.currColor = NoColor
	self.currCondition = true
}

func (self *Logger) SetColor(name string) *Logger {
	self.mu.Lock()
	self.currColor = matchColor(name)
	self.mu.Unlock()

	return self
}

func (self *Logger) Debugf(format string, v ...interface{}) {

	self.Log(Level_Debug, fmt.Sprintf(format, v...))
}

func (self *Logger) Infof(format string, v ...interface{}) {

	self.Log(Level_Info, fmt.Sprintf(format, v...))
}

func (self *Logger) Warnf(format string, v ...interface{}) {

	self.Log(Level_Warn, fmt.Sprintf(format, v...))
}

func (self *Logger) Errorf(format string, v ...interface{}) {

	self.Log(Level_Error, fmt.Sprintf(format, v...))
}

func (self *Logger) Debugln(v ...interface{}) {

	self.Log(Level_Debug, fmt.Sprintln(v...))
}

func (self *Logger) Infoln(v ...interface{}) {

	self.Log(Level_Info, fmt.Sprintln(v...))
}

func (self *Logger) Warnln(v ...interface{}) {

	self.Log(Level_Warn, fmt.Sprintln(v...))
}

func (self *Logger) Errorln(v ...interface{}) {
	self.Log(Level_Error, fmt.Sprintln(v...))
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

// 注意, 加色只能在Gogland的main方式启用, Test方式无法加色
func (self *Logger) SetColorFile(file *ColorFile) {
	self.colorFile = file
}
func (self *Logger) IsDebugEnabled() bool {
	return self.level == Level_Debug
}
