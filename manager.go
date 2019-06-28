package golog

import (
	"regexp"
	"sync"
)

var (
	loggerByName    sync.Map // map[string]*Logger
	loggerByPkgName sync.Map
)

func add(l *Logger) {

	if l.name != "" {

		if _, ok := loggerByName.Load(l.name); ok {
			panic("duplicate logger name:" + l.name)
		}

		loggerByName.Store(l.name, l)
		loggerByPkgName.Store(l.pkgName, l)
	}
}

func LoggerByName(name string) *Logger {

	if raw, ok := loggerByName.Load(name); ok {
		return raw.(*Logger)
	}

	return nil
}

// 支持正则表达式查找logger， a|b|c指定多个日志, .表示所有日志
func VisitLogger(names string, callback func(*Logger) bool) error {

	exp, err := regexp.Compile(names)
	if err != nil {
		return err
	}

	var ret []*Logger

	loggerByName.Range(func(key, value interface{}) bool {

		l := value.(*Logger)

		if exp.MatchString(l.Name()) {
			ret = append(ret, l)
		}

		return true
	})

	for _, l := range ret {
		if !callback(l) {
			break
		}
	}

	return nil
}

func ClearAll() {

	loggerByName = sync.Map{}
}
