package golog

func logPart_Text(log *Logger) {

	log.WriteRawString(log.currText)
}

func logPart_Line(log *Logger) {

	l := len(log.currText)

	if (l > 0 && log.currText[l-1] != '\n') || l == 0 {
		log.WriteRawByte('\n')
	}

}
