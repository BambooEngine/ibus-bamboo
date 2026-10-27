package introspectors

/*
#cgo CFLAGS: -std=gnu99
#cgo LDFLAGS: -lX11 -pthread
#include <stdlib.h>

extern void setXIgnoreErrorHandler();
extern char* x11GetFocusWindowClass();
extern void x11StartWindowInspector();
extern void x11StopWindowInspector();
*/
import "C"

func SetXIgnoreErrorHandler() {
	C.setXIgnoreErrorHandler()
}

func X11StartWindowInspector() {
	C.x11StartWindowInspector()
}

func X11StopWindowInspector() {
	C.x11StopWindowInspector()
}

func X11GetFocusWindowClass() string {
	var wmClass = C.x11GetFocusWindowClass()
	if wmClass != nil {
		return C.GoString(wmClass)
	}
	return ""
}
