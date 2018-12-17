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
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lX11 -lXtst
#include <sys/time.h>
#include <time.h>
#include <stdio.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/keysym.h> //xproto-devel
#include <X11/Intrinsic.h>
#include <X11/extensions/XTest.h>
inline void ucharfree(unsigned char* uc) {
	XFree(uc);
}

inline void windowfree(Window* w) {
	XFree(w);
}

inline char* uchar2char(unsigned char* uc, unsigned long len) {
	for (int i=0; i<len; i++) {
		if (uc[i] == 0) {
			uc[i] = '\n';
		}
	}
	return (char*)uc;
}

inline unsigned long uchar2long(unsigned char* uc) {
	return *(unsigned long*)(uc);
}

static int ignore_x_error(Display *display, XErrorEvent *error) {
    return 0;
}

void setXIgnoreErrorHandler() {
	XSetErrorHandler(ignore_x_error);
}

static void delay(int sec, long msec) {
    long pause;
    clock_t now,then;

    pause = msec*(CLOCKS_PER_SEC/1000);
    now = then = clock();
    while( (now-then) < pause )
        now = clock();
}

void SendBackspace(int n) {
	Display *display = XOpenDisplay(NULL);
	if (display) {
		//find out window with current focus:
        Window winfocus;
        KeyCode modcode;
        int    revert;
        XGetInputFocus(display, &winfocus, &revert);

        modcode = XKeysymToKeycode(display, XStringToKeysym("BackSpace"));
        for (int i=0; i<n; i++) {
            XTestFakeKeyEvent(display, modcode, True, 0);
            XTestFakeKeyEvent(display, modcode, False, 0);
        }
		XSync(display, 1);
		//delay(0, 1000);

		XCloseDisplay(display);
	}
}


void SendText(char* str, int len) {
	Display *display = XOpenDisplay(NULL);
	if (display) {
		//find out window with current focus:
        Window winfocus;
        int    revert;
        XGetInputFocus(display, &winfocus, &revert);

		for (int i=0; i<len; i++) {
			int modcode = str[i];
			XTestFakeKeyEvent(display, modcode, True, 0);
			XTestFakeKeyEvent(display, modcode, False, 0);
		}
		XSync(display, 1);
		//delay(0, 1000);

		XCloseDisplay(display);
	}
}
*/
import "C"
import (
	"strings"
	"unsafe"
)

const (
	MaxPropertyLen = 128
	MaxCacheWM     = 16

	WM_CLASS = "WM_CLASS"
)

type CDisplay *C.Display

var cacheWM = NewCacheWM(MaxCacheWM)

func init() {
	C.setXIgnoreErrorHandler()
}

func x11Sync(display *C.Display) {
	C.XSync(display, 0)
}

func x11Flush(display *C.Display) {
	C.XFlush(display)
}

func x11SendBackspace(n uint32) {
	C.SendBackspace(C.int(n))
}

func x11SendText(str string) {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	l := len([]rune(str))
	C.SendText(cs, C.int(l))
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
