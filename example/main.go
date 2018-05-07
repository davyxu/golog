package main

import "github.com/davyxu/golog"

var log = golog.New("test2")

// Goland中，只能在非Test环境才能使用颜色
func main() {

	golog.SetColorFile(".", "../color_sample.json")

	// 默认颜色是关闭的
	log.SetParts()
	golog.EnableColorLogger(".", true)
	log.Debugln("关闭所有部分样式")

	log.SetParts(golog.LogPart_Level)
	log.SetColor("blue")
	log.Debugln("蓝色的字+级别")

	log.SetParts(golog.LogPart_Level, golog.LogPart_Name)
	// 颜色只会影响一行
	log.SetColor("red")
	log.Warnf("级别颜色高于手动设置 + 日志名字")

	log.SetParts(golog.LogPart_Level, golog.LogPart_Name, golog.LogPart_Time, golog.LogPart_ShortFileName)
	log.Debugln()
	log.Debugf("[DB] DB日志是绿色的，从文件读取，按文字匹配的， 完整的日志样式")

	log.SetParts(golog.LogPart_TimeMS, golog.LogPart_LongFileName, func(l *golog.Logger) {
		l.WriteRawString("固定头部: ")
	})

	log.SetColor("purple")

	log.Debugf("自定义紫色 + 固定头部内容")
	log.Debugf("自定义紫色 + 固定头部内容2")

}
