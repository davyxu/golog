package golog

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFileCreate(t *testing.T) {

	l := New("test")
	SetOutputToFile("tt.log", OutputFileOption{
		MaxFileSize: 1000,
	})

	for i := 0; i < 100; i++ {
		l.Debugln(strings.Repeat(strconv.Itoa(i), 100))
	}
}
