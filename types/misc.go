package types

import (
	"sync"
)

type MacroTable struct {
	sync.RWMutex
	Enable              bool
	AutoCapitalizeMacro bool
	MTable              map[string]string
}

type TrieNode struct {
	IsWord   bool
	Value    string
	Children map[rune]*TrieNode
}
