package golog

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(log *Logger, i int, wid int) {
	var u = uint(i)
	if u == 0 && wid <= 1 {
		log.WriteRawByte('0')
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
	log.buf = append(log.buf, b[bp:]...)
}
