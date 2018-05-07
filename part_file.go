package golog

import "runtime"

func LogPart_ShortFileName(log *Logger) {

	writeFilePart(log, true, false)
}

func LogPart_LongFileName(log *Logger) {

	writeFilePart(log, false, true)
}

func writeFilePart(log *Logger, shortFile, longFile bool) {
	if shortFile || longFile {

		var file string
		var line int

		if shortFile || longFile {
			// release lock while getting caller info - it'text expensive.

			var ok bool
			_, file, line, ok = runtime.Caller(4)
			if !ok {
				file = "???"
				line = 0
			}
		}

		if shortFile {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		log.WriteRawString(file)
		log.WriteRawByte(':')
		itoa(log, line, -1)
		log.WriteRawString(": ")
	}
}
