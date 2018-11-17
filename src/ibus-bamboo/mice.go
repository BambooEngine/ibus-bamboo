package main

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lX11
extern void mouse_capture_init();
extern void mouse_capture_exit();
extern void mouse_capture_unlock();
*/
import "C"

//export mouse_move_handler
func mouse_move_handler() {
	go onMouseMove()
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
