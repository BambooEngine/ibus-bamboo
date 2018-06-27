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
	"fmt"
	"os"
	"strings"
)

const IBusDaemon = "ibus-daemon"

func getProcessName(pid int) string {
	pStatusFile := fmt.Sprintf("/proc/%d/status", pid)
	f, e := os.OpenFile(pStatusFile, os.O_RDONLY, os.ModePerm)
	if e == nil {
		defer f.Close()
		buf := make([]byte, len(IBusDaemon)*2)
		n, e := f.Read(buf)
		if e == nil {
			s := string(buf[:n])
			lines := strings.Split(s, "\n")
			firstLineParts := strings.Split(lines[0], "\t")
			if len(firstLineParts) >= 2 {
				return firstLineParts[1]
			}
		}
	}

	return ""
}

func isIBusDaemonChild() bool {
	ppid := os.Getppid()
	return getProcessName(ppid) == IBusDaemon
}
