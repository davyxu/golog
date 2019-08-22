package golog

func LogPart_CurrLevel(log *Logger) {
	log.WriteRawString(log.CurrLevelString())
	log.WriteRawByte(' ')

}

func LogPart_Name(log *Logger) {

	if log.name != "" {
		log.WriteRawString(log.name)
		log.WriteRawByte(' ')
	}
}
