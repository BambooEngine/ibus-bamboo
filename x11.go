/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

/*
#cgo CFLAGS: -std=gnu99
#cgo LDFLAGS: -lX11 -lXtst -pthread
#include <stdlib.h>

extern void x11Copy(char*);
extern void x11Paste(int);
extern void clipboard_init();
extern void clipboard_exit();
extern void x11ClipboardReset();
extern void x11SendBackspace(int n, int timeout);
extern void x11SendShiftR();
extern void x11SendShiftLeft(int n, int r, int timeout);
*/
import "C"
import (
	"ibus-bamboo/introspectors"
	"sync"
	"unsafe"
)

var (
	mcapRunning bool
	mcapMutex   sync.RWMutex
)

func init() {
	introspectors.SetXIgnoreErrorHandler()
}

func x11StartWindowInspector() {
	introspectors.X11StartWindowInspector()
}

func x11StopWindowInspector() {
	introspectors.X11StopWindowInspector()
}

func x11Copy(str string) {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	C.x11Copy(cs)
}

func x11ClipboardInit() {
	C.clipboard_init()
}

func x11ClipboardExit() {
	C.clipboard_exit()
}

func x11ClipboardReset() {
	C.x11ClipboardReset()
}

func x11Paste(n int) {
	C.x11Paste(C.int(n))
}

func x11SendShiftR() {
	C.x11SendShiftR()
}

func x11SendShiftLeft(n int, shiftRightIsPressing bool, timeout int) {
	var r = 0
	if shiftRightIsPressing {
		r = 1
	}
	C.x11SendShiftLeft(C.int(n), C.int(r), C.int(timeout))
}

func x11SendBackspace(n int, timeout int) {
	C.x11SendBackspace(C.int(n), C.int(timeout))
}

func x11GetFocusWindowClass() string {
	return introspectors.X11GetFocusWindowClass()
}
