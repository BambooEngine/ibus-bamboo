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

var wordBreakSyms = []rune{
	',', ';', ':', '.', '"', '\'', '!', '?', ' ',
	'<', '>', '=', '+', '-', '*', '/', '\\',
	'_', '~', '`', '@', '#', '$', '%', '^', '&', '(', ')', '{', '}', '[', ']',
	'|',
}

func isWordBreakSymbol(key rune) bool {
	if key >= '0' && key <= '9' {
		return true
	}
	for _, c := range wordBreakSyms {
		if c == key {
			return true
		}
	}
	return false
}

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

func inWhiteList(list []string, classes string) bool {
	for _, cl := range list {
		if cl == classes {
			return true
		}
	}
	return false
}

func isSameClasses(cl1 string, cl2 string) bool {
	return cl1 == cl2
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
