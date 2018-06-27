/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Nguyen Cong Hoang <hoangnc.jp@gmail.com>
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

import (
	"os"
)

//sudo chmod +r /dev/input/mice
const (
	DevInputMice = "/dev/input/mice"
)

var onMouseClick func()
var onMouseMove func()

func init() {
	go func() {
		down := false
		miceDev, err := os.OpenFile(DevInputMice, os.O_RDONLY, 0)
		if err == nil {
			data := make([]byte, 3)
			for {
				n, err := miceDev.Read(data)
				if err == nil && n == 3 {
					go onMouseMove()
				}
				if err == nil && n == 3 && data[0]&0x7 != 0 {
					if data[1] == 0 && data[2] == 0 {
						if !down {
							if onMouseClick != nil {
								go onMouseClick()
							}
							down = true
						}
					}
				} else if down {
					down = false
				}
			}
		}
	}()
}
