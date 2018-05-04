package golog

import (
	"runtime"
	"time"
)

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *[]byte, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		*buf = append(*buf, '0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	*buf = append(*buf, b[bp:]...)
}

func writeTimePart(flag LogFlag, buf *[]byte) {

	if flag.Contains(Ldate | Ltime | Lmicroseconds) {

		now := time.Now() // get this early.

		if flag.Contains(Ldate) {
			year, month, day := now.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if flag.Contains(Ltime | Lmicroseconds) {
			hour, min, sec := now.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if flag.Contains(Lmicroseconds) {
				*buf = append(*buf, '.')
				itoa(buf, now.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
}

func writeFilePart(flag LogFlag, buf *[]byte) {
	if flag.Contains(Lshortfile | Llongfile) {

		var file string
		var line int

		if flag.Contains(Lshortfile | Llongfile) {
			// release lock while getting caller info - it'text expensive.

			var ok bool
			_, file, line, ok = runtime.Caller(4)
			if !ok {
				file = "???"
				line = 0
			}
		}

		if flag.Contains(Lshortfile) {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}
