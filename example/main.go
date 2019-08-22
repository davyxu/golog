package main

import "github.com/davyxu/golog"

var log = golog.New("test2")

// Goland中，只能在非Test环境才能使用颜色

const colorStyle = `
{
	"Rule":[
		{"Text":"panic:","Color":"Red"},
		{"Text":"[DB]","Color":"Green"},
		{"Text":"#http.listen","Color":"Blue"},
		{"Text":"#http.recv","Color":"Blue"},
		{"Text":"#http.send","Color":"Purple"},

		{"Text":"#tcp.listen","Color":"Blue"},
		{"Text":"#tcp.accepted","Color":"Blue"},
		{"Text":"#tcp.closed","Color":"Blue"},
		{"Text":"#tcp.recv","Color":"Blue"},
		{"Text":"#tcp.send","Color":"Purple"},
		{"Text":"#tcp.connected","Color":"Blue"},

		{"Text":"#udp.listen","Color":"Blue"},
		{"Text":"#udp.recv","Color":"Blue"},
		{"Text":"#udp.send","Color":"Purple"},

		{"Text":"#rpc.recv","Color":"Blue"},
		{"Text":"#rpc.send","Color":"Purple"},

		{"Text":"#relay.recv","Color":"Blue"},
		{"Text":"#relay.send","Color":"Purple"}
	]
}
`

func main() {

	golog.SetColorDefine(".", colorStyle)

	// 默认颜色是关闭的
	log.SetParts()
	golog.EnableColorLogger(".", true)
	log.Debugln("关闭所有部分样式")

	log.SetParts(golog.LogPart_CurrLevel)
	log.SetColor("blue")
	log.Debugln("蓝色的字+级别")

	log.SetParts(golog.LogPart_CurrLevel, golog.LogPart_Name)
	// 颜色只会影响一行
	log.SetColor("red")
	log.Warnf("级别颜色高于手动设置 + 日志名字")

	log.SetParts(golog.LogPart_CurrLevel, golog.LogPart_Name, golog.LogPart_Time, golog.LogPart_ShortFileName)
	log.Debugln()
	log.Debugf("[DB] DB日志是绿色的，从文件读取，按文字匹配的， 完整的日志样式")

	log.SetParts(golog.LogPart_TimeMS, golog.LogPart_LongFileName, func(l *golog.Logger) {
		l.WriteRawString("固定头部: ")
	})

	log.SetColor("purple")

	log.Debugf("自定义紫色 + 固定头部内容")
	log.Debugf("自定义紫色 + 固定头部内容2")

}
