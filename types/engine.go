package types

import (
	"sync"

	"github.com/BambooEngine/bamboo-core"
	ibus "github.com/BambooEngine/goibus"
	"github.com/godbus/dbus/v5"
)

type IBusBambooEngineState struct {
	sync.Mutex
	IEngine
	Preeditor              bamboo.IEngine
	EngineName             string
	Config                 *Config
	PropList               *ibus.PropList
	EnglishMode            bool
	MacroTable             *MacroTable
	WmClasses              string
	IsInputModeLTOpened    bool
	IsEmojiLTOpened        bool
	IsInHexadecimal        bool
	EmojiLookupTable       *ibus.LookupTable
	InputModeLookupTable   *ibus.LookupTable
	Capabilities           uint32
	KeyPressDelay          int
	NFakeBackSpace         int32
	IsFirstTimeSendingBS   bool
	Emoji                  *EmojiEngine
	IsSurroundingTextReady bool
	LastKeyWithShift       bool
	LastCommitText         int64
	// restore key strokes by pressing Shift + Space
	ShouldRestoreKeyStrokes bool
	// enqueue key strokes to process later
	ShouldEnqueuKeyStrokes bool
}

type IEngine interface {
	GetAll(iface string) (map[string]dbus.Variant, *dbus.Error)
	ProcessKeyEvent(keyval uint32, keycode uint32, state uint32) (bool, *dbus.Error)
	SetCursorLocation(x int32, y int32, w int32, h int32) *dbus.Error
	SetSurroundingText(text dbus.Variant, cursor_index uint32, anchor_pos uint32) *dbus.Error
	SetCapabilities(cap uint32) *dbus.Error
	FocusIn() *dbus.Error
	FocusOut() *dbus.Error
	Reset() *dbus.Error
	PageUp() *dbus.Error
	PageDown() *dbus.Error
	CursorUp() *dbus.Error
	CursorDown() *dbus.Error
	CandidateClicked(index uint32, button uint32, state uint32) *dbus.Error
	Enable() *dbus.Error
	Disable() *dbus.Error
	PropertyActivate(prop_name string, prop_state uint32) *dbus.Error
	PropertyShow(prop_name string) *dbus.Error
	PropertyHide(prop_name string) *dbus.Error
	Destroy() *dbus.Error
	CommitText(text *ibus.Text)
	ForwardKeyEvent(keyval uint32, keycode uint32, state uint32)
	UpdatePreeditText(text *ibus.Text, cursor_pos uint32, visible bool)
	UpdatePreeditTextWithMode(text *ibus.Text, cursor_pos uint32, visible bool, mode uint32)
	ShowPreeditText()
	HidePreeditText()
	UpdateAuxiliaryText(text *ibus.Text, visible bool)
	ShowAuxiliaryText()
	HideAuxiliaryText()
	UpdateLookupTable(lookup_table *ibus.LookupTable, visible bool)
	ShowLookupTable()
	HideLookupTable()
	PageUpLookupTable()
	PageDownLookupTable()
	CursorUpLookupTable()
	CursorDownLookupTable()
	RegisterProperties(props *ibus.PropList)
	UpdateProperty(prop *ibus.Property)
	DeleteSurroundingText(offset_from_cursor int32, nchars uint32)
	RequireSurroundingText()
}
