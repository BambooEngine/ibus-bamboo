package types

import (
	"github.com/BambooEngine/bamboo-core"
)

type Config struct {
	InputMethod            string
	InputMethodDefinitions map[string]bamboo.InputMethodDefinition
	OutputCharset          string
	Flags                  uint
	IBflags                uint
	Shortcuts              [10]uint32
	DefaultInputMode       int
	InputModeMapping       map[string]int
	Introspector           BambooIntrospector
	IntrospectorCommand    string
}
