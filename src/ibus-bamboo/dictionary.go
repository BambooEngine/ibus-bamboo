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

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
)

//var rootWordTrie = &bamboo.W{F: false}

func fileExist(p string) bool {
	sta, err := os.Stat(p)
	return err == nil && !sta.IsDir()
}

func loadDictionary(dataFiles ...string) (map[string]bool, error) {
	var dictionary = map[string]bool{}
	for _, dataFile := range dataFiles {
		if !fileExist(dataFile) && !filepath.IsAbs(dataFile) {
			dataFile = filepath.Join(filepath.Dir(os.Args[0]), dataFile)
		}
		f, err := os.Open(dataFile)
		if err != nil {
			return nil, err
		}
		rd := bufio.NewReader(f)
		for {
			line, _, err := rd.ReadLine()
			if err != nil {
				break
			}
			if len(line) == 0 {
				continue
			}
			dictionary[string(line)] = true
			//bamboo.AddTrie(rootWordTrie, []rune(string(line)), false)
		}
		f.Close()
	}
	runtime.GC()
	debug.FreeOSMemory()
	return dictionary, nil
}
