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
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type EmojiOne struct {
	Name      string
	Shortname string
	Keywords  []string
	Ascii     []string
}

func loadEmojiOne(dataFile string) {
	emojiTrie = NewTrie()
	var c = map[string]EmojiOne{}
	var data, err = ioutil.ReadFile(dataFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &c)
	for k, v := range c {
		var codePointStr string
		for _, codePoint := range strings.Split(k, "-") {
			if code, err := strconv.ParseInt(codePoint, 16, 32); err == nil {
				codePointStr += string(rune(code))
			}
		}
		for _, ascii := range v.Ascii {
			InsertTrie(emojiTrie, ascii, codePointStr)
		}
		for _, keyword := range v.Keywords {
			InsertTrie(emojiTrie, keyword, codePointStr)
		}
	}
}

type EmojiEngine struct {
	keys []rune
}

func NewEmojiEngine() *EmojiEngine {
	var be = &EmojiEngine{}
	return be
}

func (be *EmojiEngine) MatchString(s string) bool {
	var lookup = FindPrefix(emojiTrie, s)
	return lookup != nil
}

func (be *EmojiEngine) Filter(s string) []string {
	var codePoints []string
	var keys []string
	var lookup = FindPrefix(emojiTrie, s)
	for key, _ := range lookup {
		keys = append(keys, key)
	}
	var names = byString(keys)
	sort.Sort(names)
	for _, name := range names {
		codePoints = append(codePoints, lookup[name])
	}
	return codePoints
}

func (be *EmojiEngine) ProcessKey(key rune) {
	be.keys = append(be.keys, key)
}

func (be *EmojiEngine) GetRawString() string {
	return string(be.keys)
}

func (be *EmojiEngine) Reset() {
	be.keys = nil
}

func (be *EmojiEngine) Query() []string {
	return be.Filter(string(be.keys))
}

func (be *EmojiEngine) RemoveLastKey() {
	if len(be.keys) <= 0 {
		return
	}
	be.keys = be.keys[:len(be.keys)-1]
}
