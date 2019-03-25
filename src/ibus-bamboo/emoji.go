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
	"os"
	"path/filepath"
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

func loadEmojiOne(dataFile string) (map[string]EmojiOne, error) {
	var c = map[string]EmojiOne{}
	if !fileExist(dataFile) && !filepath.IsAbs(dataFile) {
		dataFile = filepath.Join(filepath.Dir(os.Args[0]), dataFile)
	}
	if data, err := ioutil.ReadFile(dataFile); err == nil {
		json.Unmarshal(data, &c)
	}
	return c, nil
}

type BambooEmoji struct {
	nameTable      map[string]string
	shortNameTable map[string]string
	asciiTable     map[string]string
	emojiTrie      *bamboo.W
	keys           []rune
}

func NewBambooEmoji(emojiDictPath string) *BambooEmoji {
	var be = &BambooEmoji{
		shortNameTable: map[string]string{},
		asciiTable:     map[string]string{},
		emojiTrie:      &bamboo.W{},
	}
	var data, err = loadEmojiOne(emojiDictPath)
	if err != nil {
		return be
	}
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
			bamboo.AddTrie(be.emojiTrie, []rune(ascii), false)
		}
		bamboo.AddTrie(be.emojiTrie, []rune(shortName), false)
	}
	return be
}

func (be *BambooEmoji) TestString(s string) uint8 {
	return bamboo.TestString(be.emojiTrie, []rune(s), false)
}

func (be *BambooEmoji) Filter(s string) []string {
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

func (be *BambooEmoji) ProcessKey(key rune) {
	be.keys = append(be.keys, key)
}

func (be *BambooEmoji) GetRawString() string {
	return string(be.keys)
}

func (be *BambooEmoji) Reset() {
	be.keys = nil
}

func (be *BambooEmoji) Query() []string {
	return be.Filter(string(be.keys))
}

func (be *BambooEmoji) RemoveLastKey() {
	if len(be.keys) <= 0 {
		return
	}
	be.keys = be.keys[:len(be.keys)-1]
}
