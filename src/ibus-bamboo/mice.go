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
#cgo LDFLAGS: -lX11
extern void mouse_capture_init();
extern void mouse_capture_exit();
extern void mouse_capture_unlock();
*/
import "C"

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
