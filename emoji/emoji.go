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

package emoji

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"ibus-bamboo/types"
)

var Trie = NewTrie()

func LoadEmojiOne(dataFile string) (*types.TrieNode, error) {
	var trie = NewTrie()
	var c = map[string]types.EmojiOne{}
	var data, err = ioutil.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &c)
	for k, v := range c {
		var codePointStr string
		for _, codePoint := range strings.Split(k, "-") {
			if code, err := strconv.ParseInt(codePoint, 16, 32); err == nil {
				codePointStr += string(rune(code))
			}
		}
		for _, ascii := range v.ASCII {
			InsertTrie(trie, ascii, codePointStr)
		}
		for _, keyword := range v.Keywords {
			InsertTrie(trie, keyword, codePointStr)
		}
	}
	return trie, nil
}

func NewEmojiEngine() *types.EmojiEngine {
	var be = &types.EmojiEngine{}
	return be
}

func MatchString(be *types.EmojiEngine, s string) bool {
	var lookup = FindPrefix(Trie, s)
	return lookup != nil
}

func Filter(be *types.EmojiEngine, s string) []string {
	var codePoints []string
	var keys []string
	var lookup = FindPrefix(Trie, s)
	for key := range lookup {
		keys = append(keys, key)
	}
	var names = sort.StringSlice(keys)
	sort.Sort(names)
	for _, name := range names {
		var cps = sort.StringSlice(strings.Split(lookup[name], ":"))
		sort.Sort(cps)
		for _, cp := range cps {
			codePoints = append(codePoints, cp)
		}
	}
	return codePoints
}

func ProcessKey(be *types.EmojiEngine, key rune) {
	be.Keys = append(be.Keys, key)
}

func GetRawString(be *types.EmojiEngine) string {
	return string(be.Keys)
}

func Reset(be *types.EmojiEngine) {
	be.Keys = nil
}

func Query(be *types.EmojiEngine) []string {
	return Filter(be, string(be.Keys))
}

func RemoveLastKey(be *types.EmojiEngine) {
	if len(be.Keys) <= 0 {
		return
	}
	be.Keys = be.Keys[:len(be.Keys)-1]
}
