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

type TrieNode struct {
	isWord   bool
	value    string
	Children map[rune]*TrieNode
}

func NewTrie() *TrieNode {
	var trie = new(TrieNode)
	trie.Children = make(map[rune]*TrieNode)
	return trie
}

func InsertTrie(r *TrieNode, word, value string) {
	var currentNode = r
	for _, c := range word {
		if cn := currentNode.Children[c]; cn != nil {
			currentNode = cn
		} else {
			var n = new(TrieNode)
			n.Children = make(map[rune]*TrieNode)
			currentNode.Children[c] = n
			currentNode = n
		}
	}
	if currentNode.value == "" {
		currentNode.value = value
	} else {
		currentNode.value += ":" + value
	}
	currentNode.isWord = true
}

func dfs(trie *TrieNode, lookup map[string]string, s string) {
	if trie.isWord {
		lookup[s] = trie.value
	}
	for chr, t := range trie.Children {
		var key = s + string(chr)
		dfs(t, lookup, key)
	}
}

func FindPrefix(r *TrieNode, prefix string) map[string]string {
	var currentNode = r
	for _, c := range prefix {
		if cn := currentNode.Children[c]; cn != nil {
			currentNode = cn
		} else {
			return nil
		}
	}
	var lookup = make(map[string]string)
	dfs(currentNode, lookup, prefix)
	return lookup
}
