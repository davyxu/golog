package golog

func LogPart_Level(log *Logger) {
	log.WriteRawString(levelString[log.currLevel])
	log.WriteRawByte(' ')

}

func LogPart_Name(log *Logger) {

	if log.name != "" {
		log.WriteRawString(log.name)
		log.WriteRawByte(' ')
	}
}
