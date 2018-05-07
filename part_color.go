package golog

func logPart_ColorBegin(log *Logger) {

	if log.enableColor && log.currColor != NoColor {

		log.WriteRawString(logColorPrefix[log.currColor])
	}
}

func logPart_ColorEnd(log *Logger) {

	if log.enableColor && log.currColor != NoColor {

		log.WriteRawString(logColorSuffix)
	}
}
