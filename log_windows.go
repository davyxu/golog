package golog

import (
	"errors"
)

/*
#include <stdio.h>
#include <windows.h>

int EnableVT100() {
	const DWORD ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004;

	HANDLE hOut = GetStdHandle(STD_OUTPUT_HANDLE);
    if (hOut == INVALID_HANDLE_VALUE) {
        return 1;
    }

    DWORD dwMode = 0;
    if (!GetConsoleMode(hOut, &dwMode)) {
        return 2;
    }

    dwMode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING;
    if (!SetConsoleMode(hOut, dwMode)) {
        return 3;
    }

	return 0;
}
*/
import "C"

func init() {
	EnableVT100()
}

func EnableVT100() error {
	r := int(C.EnableVT100())
	if r != 0 {
		return errors.New("enable VT100 support failed")
	}
	return nil
}
