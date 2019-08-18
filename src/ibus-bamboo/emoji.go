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
	"github.com/BambooEngine/bamboo-core"
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

var emojiMap map[string]EmojiOne

func loadEmojiOne(dataFile string) (map[string]EmojiOne, error) {
	var c = map[string]EmojiOne{}
	if data, err := ioutil.ReadFile(dataFile); err == nil {
		json.Unmarshal(data, &c)
	}
	return c, nil
}

type EmojiEngine struct {
	nameTable      map[string]string
	shortNameTable map[string]string
	asciiTable     map[string]string
	emojiTrie      *bamboo.Node
	keys           []rune
}

func NewEmojiEngine() *EmojiEngine {
	var be = &EmojiEngine{
		shortNameTable: map[string]string{},
		asciiTable:     map[string]string{},
		emojiTrie:      &bamboo.Node{},
	}
	var data = emojiMap
	for k, v := range data {
		var codePointStr string
		for _, codePoint := range strings.Split(k, "-") {
			if code, err := strconv.ParseInt(codePoint, 16, 32); err == nil {
				codePointStr += string(rune(code))
			}
		}
		var shortName = v.Shortname[1 : len([]rune(v.Shortname))-1]
		be.shortNameTable[shortName] = codePointStr
		for _, ascii := range v.Ascii {
			be.asciiTable[ascii] = codePointStr
			bamboo.AddTrie(be.emojiTrie, []rune(ascii), false, false)
		}
		bamboo.AddTrie(be.emojiTrie, []rune(shortName), false, false)
	}
	return be
}

func (be *EmojiEngine) TestString(s string) uint8 {
	return bamboo.TestString(be.emojiTrie, []rune(s), false)
}

func (be *EmojiEngine) Filter(s string) []string {
	var codePoints []string
	var names = byString(bamboo.FindWords(be.emojiTrie, s))
	sort.Sort(names)
	for _, name := range names {
		if be.asciiTable[name] != "" {
			codePoints = append(codePoints, be.asciiTable[name])
		} else if be.shortNameTable[name] != "" {
			codePoints = append(codePoints, be.shortNameTable[name])
		}
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
