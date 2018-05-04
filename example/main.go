package main

import "github.com/davyxu/golog"

var log = golog.New("test2")

// Goland中，只能在非Test环境才能使用颜色
func main() {

	golog.SetColorFile(".", "../color_sample.json")

	// 默认颜色是关闭的
	golog.EnableColorLogger(".", true)

	log.SetFlag(0)

	log.Debugln("关闭所有部分样式")

	log.SetFlag(golog.Llevel)
	log.SetColor("blue")
	log.Debugln("蓝色的字+级别")

	log.SetFlag(golog.Llevel | golog.Lname)
	// 颜色只会影响一行
	log.SetColor("red")
	log.Warnf("级别颜色高于手动设置 + 日志名字")

	log.SetFlag(golog.LstdFlags | golog.Lshortfile)
	log.Debugf("[DB] DB日志是绿色的，从文件读取，按文字匹配的， 完整的日志样式")

}
