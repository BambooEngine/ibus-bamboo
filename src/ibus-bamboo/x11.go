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
#include <X11/Xlib.h>

extern void x11Copy(char*);
extern void clipboard_init();
extern void clipboard_exit();
extern void mouse_capture_init();
extern void mouse_capture_exit();
extern void mouse_capture_unlock();
extern unsigned long uchar2long(unsigned char* uc);
extern char* uchar2char(unsigned char* uc, unsigned long len);
extern void windowfree(Window* w);
extern void ucharfree(unsigned char* uc);
extern void send_text(char* str);
extern void x11Paste(int);
extern void send_backspace(int n);
extern void setXIgnoreErrorHandler();
*/
import "C"
import (
	"strings"
	"unsafe"
)

const (
	MaxPropertyLen = 128

	WM_CLASS = "WM_CLASS"
)

type CDisplay *C.Display

func init() {
	C.setXIgnoreErrorHandler()
}

//export mouse_move_handler
func mouse_move_handler() {
	onMouseMove()
}

var onMouseMove func()

func mouseCaptureInit() {
	C.mouse_capture_init()
}

func mouseCaptureExit() {
	C.mouse_capture_exit()
}

func mouseCaptureUnlock() {
	C.mouse_capture_unlock()
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

func x11Paste(n int) {
	C.x11Paste(C.int(n))
}

func x11SendBackspace(n uint32) {
	C.send_backspace(C.int(n))
}

func x11SendText(str string) {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	C.send_text(cs)
}

func x11GetUCharProperty(display *C.Display, window C.Window, propName string) (*C.uchar, C.ulong) {
	var actualType C.Atom
	var actualFormat C.int
	var nItems, bytesAfter C.ulong
	var prop *C.uchar

	filterAtom := C.XInternAtom(display, C.CString(propName), C.True)

	status := C.XGetWindowProperty(display, window, filterAtom, 0, MaxPropertyLen, C.False, C.AnyPropertyType, &actualType, &actualFormat, &nItems, &bytesAfter, &prop)

	if status == C.Success {
		return prop, nItems
	}

	return nil, 0
}

func x11GetStringProperty(display *C.Display, window C.Window, propName string) string {
	prop, propLen := x11GetUCharProperty(display, window, propName)
	if prop != nil {
		defer C.ucharfree(prop)
		return C.GoString(C.uchar2char(prop, propLen))
	}

	return ""
}

func x11OpenDisplay() *C.Display {
	return C.XOpenDisplay(nil)
}

func x11GetInputFocus(display *C.Display) C.Window {
	var window C.Window
	var revertTo C.int
	C.XGetInputFocus(display, &window, &revertTo)

	return window
}

func x11GetParentWindow(display *C.Display, w C.Window) (rootWindow, parentWindow C.Window) {
	var childrenWindows *C.Window
	var nChild C.uint
	C.XQueryTree(display, w, &rootWindow, &parentWindow, &childrenWindows, &nChild)
	C.windowfree(childrenWindows)

	return
}

func x11CloseDisplay(d *C.Display) {
	C.XCloseDisplay(d)
}

func x11GetFocusWindowClasses(display *C.Display) []string {

	if display != nil {

		w := x11GetInputFocus(display)
		strClass := ""
		for {
			s := x11GetStringProperty(display, w, WM_CLASS)
			if len(s) > 0 {
				strClass += s + "\n"
			}

			rootWindow, parentWindow := x11GetParentWindow(display, w)

			if rootWindow == parentWindow {
				break
			}

			w = parentWindow
		}

		return strings.Split(strClass, "\n")
	}
	return nil
}

func x11GetFocusWindowClass(display *C.Display) []string {
	var wmClasses = x11GetFocusWindowClasses(display)
	if len(wmClasses) >= 2 {
		return []string{wmClasses[1]}
	}
	return nil
}
