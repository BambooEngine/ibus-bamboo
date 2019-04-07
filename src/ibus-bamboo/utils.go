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
	"github.com/BambooEngine/bamboo-core"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

func toUpper(keyRune rune) rune {
	var upperSpecialKeys = map[rune]rune{
		'[': '{',
		']': '}',
	}

	if upperSpecialKey, found := upperSpecialKeys[keyRune]; found {
		keyRune = upperSpecialKey
	} else {
		keyRune = unicode.ToUpper(keyRune)
	}
	return keyRune
}

func inStringList(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func inWMList(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		} else if re, err := regexp.Compile(s); err == nil && re.MatchString(str) {
			return true
		}
	}
	return false
}

func removeFromWhiteList(list []string, classes string) []string {
	var newList []string
	for _, cl := range list {
		if cl != classes {
			newList = append(newList, cl)
		}
	}
	return newList
}

func addToWhiteList(list []string, classes string) []string {
	for _, cl := range list {
		if cl == classes {
			return list
		}
	}
	return append(list, classes)
}

func getCharsetFromPropKey(str string) (string, bool) {
	var arr = strings.Split(str, "-")
	if len(arr) == 2 {
		return arr[1], true
	}
	return str, false
}

func isValidCharset(str string) bool {
	var charsets = bamboo.GetCharsetNames()
	for _, cs := range charsets {
		if cs == str {
			return true
		}
	}
	return false
}

type byString []string

func (s byString) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s byString) Len() int {
	return len(s)
}
func (s byString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sortStrings(list []string) []string {
	var strList = byString(list)
	sort.Sort(strList)
	return strList
}

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
	return dictionary, nil
}
