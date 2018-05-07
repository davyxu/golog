package golog

func LogPart_Level(log *Logger) {
	log.WriteRawString(levelString[log.currLevel])
	log.WriteRawByte(' ')

}

func LogPart_Name(log *Logger) {

	log.WriteRawString(log.name)
	log.WriteRawByte(' ')
}
