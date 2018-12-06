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

package bamboo

import (
	"unicode"
)

const (
	FindResultNotMatch = iota
	FindResultMatchPrefix
	FindResultMatchFull
)

//Word trie
type W struct {
	F bool        //Full word
	N map[rune]*W // Next characters
}

func AddTrie(trie *W, s []rune, down bool) {
	if trie.N == nil {
		trie.N = map[rune]*W{}
	}

	//add original char
	s0 := s[0]
	if trie.N[s0] == nil {
		trie.N[s0] = &W{}
	}

	if len(s) == 1 {
		if !trie.N[s0].F {
			trie.N[s0].F = !down
		}
	} else {
		AddTrie(trie.N[s0], s[1:], down)
	}

	//add down 1 level char
	if dmap, exist := downLvlMap[s0]; exist {
		for _, r := range dmap {
			if trie.N[r] == nil {
				trie.N[r] = &W{}
			}

			if len(s) == 1 {
				trie.N[r].F = true
			} else {
				AddTrie(trie.N[r], s[1:], true)
			}
		}
	}
}

func FindWord(trie *W, s []rune, deepSearch bool) uint8 {

	if len(s) == 0 {
		if trie.F {
			if deepSearch && len(trie.N) > 0 {
				return FindResultMatchPrefix
			}
			return FindResultMatchFull
		}
		return FindResultMatchPrefix
	}

	c := unicode.ToLower(s[0])

	if trie.N[c] != nil {
		r := FindWord(trie.N[c], s[1:], deepSearch)
		if r != FindResultNotMatch {
			return r
		}
	}

	return FindResultNotMatch
}
