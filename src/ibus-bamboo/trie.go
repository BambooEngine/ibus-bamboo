package main

type TrieNode struct {
	isWord   bool
	value    string
	children map[rune]*TrieNode
}

func NewTrie() *TrieNode {
	var trie = new(TrieNode)
	trie.children = make(map[rune]*TrieNode)
	return trie
}

func InsertTrie(r *TrieNode, word, value string) {
	var currentNode = r
	for _, c := range word {
		if cn := currentNode.children[c]; cn != nil {
			currentNode = cn
		} else {
			var n = new(TrieNode)
			n.children = make(map[rune]*TrieNode)
			currentNode.children[c] = n
			currentNode = n
		}
	}
	currentNode.value = value
	currentNode.isWord = true
}

func dfs(trie *TrieNode, lookup map[string]string, s string) {
	if trie.isWord {
		lookup[s] = trie.value
	}
	for chr, t := range trie.children {
		var key = s + string(chr)
		dfs(t, lookup, key)
	}
}

func FindPrefix(r *TrieNode, prefix string) map[string]string {
	var currentNode = r
	for _, c := range prefix {
		if cn := currentNode.children[c]; cn != nil {
			currentNode = cn
		} else {
			return nil
		}
	}
	var lookup = make(map[string]string)
	dfs(currentNode, lookup, prefix)
	return lookup
}
