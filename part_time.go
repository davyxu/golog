package golog

import "time"

func LogPart_Time(log *Logger) {

	writeTimePart(log, false)
}

func LogPart_TimeMS(log *Logger) {

	writeTimePart(log, true)
}

func writeTimePart(log *Logger, ms bool) {

	now := time.Now() // get this early.

	year, month, day := now.Date()

	itoa(log, year, 4)
	log.WriteRawByte('/')
	itoa(log, int(month), 2)
	log.WriteRawByte('/')
	itoa(log, day, 2)
	log.WriteRawByte(' ')

	hour, min, sec := now.Clock()
	itoa(log, hour, 2)
	log.WriteRawByte(':')
	itoa(log, min, 2)
	log.WriteRawByte(':')
	itoa(log, sec, 2)

	if ms {
		log.WriteRawByte('.')
		itoa(log, now.Nanosecond()/1e6, 3)
	}

	log.WriteRawByte(' ')

}
