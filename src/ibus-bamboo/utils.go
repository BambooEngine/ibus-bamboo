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
	"github.com/BambooEngine/bamboo-core"
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

func inKeyMap(keys []rune, key rune) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func inWhiteList(list map[string]bool, classes []string) bool {
	for _, cl := range classes {
		if list[cl] {
			return true
		}
	}
	return false
}

func isSameClasses(cl1 []string, cl2 []string) bool {
	if len(cl1) != len(cl2) || len(cl1) != 1 {
		return false
	}
	return cl1[0] == cl2[0]
}

func removeFromWhiteList(list map[string]bool, classes []string) map[string]bool {
	for _, cl := range classes {
		if list[cl] {
			list[cl] = false
		}
	}
	return list
}

func addToWhiteList(list map[string]bool, classes []string) map[string]bool {
	for _, cl := range classes {
		list[cl] = true
	}
	return list
}

func getCharsetFromPropKey(str string) string {
	var arr = strings.Split(str, "-")
	if len(arr) == 2 {
		return arr[1]
	}
	return str
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
