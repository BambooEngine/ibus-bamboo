/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LISENCE>.
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

func TestString(trie *W, s []rune, deepSearch bool) uint8 {

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
		r := TestString(trie.N[c], s[1:], deepSearch)
		if r != FindResultNotMatch {
			return r
		}
	}

	return FindResultNotMatch
}

func dfs(trie *W, lookup map[string]bool, s string) {
	if trie.F {
		lookup[s] = true
	}
	for chr, t := range trie.N {
		var key = s + string(chr)
		dfs(t, lookup, key)
	}
}

func FindNode(trie *W, s []rune) *W {
	if len(s) == 0 {
		return trie
	}
	c := s[0]
	if trie.N[c] != nil {
		return FindNode(trie.N[c], s[1:])
	}
	// not match
	return nil
}

func FindWords(trie *W, s string) []string {
	var words []string
	var node = FindNode(trie, []rune(s))
	if node == nil {
		return nil
	}
	var lookup = map[string]bool{}
	dfs(node, lookup, s)
	for w := range lookup {
		words = append(words, w)
	}
	return words
}
