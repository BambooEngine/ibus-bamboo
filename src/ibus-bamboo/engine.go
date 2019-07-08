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
	"fmt"
	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
	"log"
	"os/exec"
	"reflect"
	"sync"
)

type IBusBambooEngine struct {
	sync.Mutex
	ibus.Engine
	preeditor            bamboo.IEngine
	zeroLocation         bool
	engineName           string
	config               *Config
	propList             *ibus.PropList
	mode                 bamboo.Mode
	ignorePreedit        bool
	macroTable           *MacroTable
	dictionary           map[string]bool
	wmClasses            string
	isInputModeLTOpened  bool
	isEmojiLTOpened      bool
	emojiLookupTable     *ibus.LookupTable
	inputModeLookupTable *ibus.LookupTable
	capabilities         uint32
	nFakeBackSpace       int
	firstTimeSendingBS   bool
	isFocusOut           bool
	emoji                *BambooEmoji
	lastKeyWithShift     bool
}

/**
Implement IBus.Engine's process_key_event default signal handler.

Args:
	keyval - The keycode, transformed through a keymap, stays the
		same for every keyboard
	keycode - Keyboard-dependant key code
	modifiers - The state of IBus.ModifierType keys like
		Shift, Control, etc.
Return:
	True - if successfully process the keyevent
	False - otherwise. The keyevent will be passed to X-Client

This function gets called whenever a key is pressed.
*/
func (e *IBusBambooEngine) ProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	if e.isIgnoredKey(keyVal, state) {
		return false, nil
	}
	log.Printf("keyCode 0x%04x keyval 0x%04x | %c | %d\n", keyCode, keyVal, rune(keyVal), len(keyPressChan))
	if e.config.IBflags&IBinputLookupTableDisabled != 0 && keyVal == IBUS_OpenLookupTable && e.isInputModeLTOpened == false && e.wmClasses != "" {
		e.resetBuffer()
		e.isInputModeLTOpened = true
		e.openLookupTable()
		return true, nil
	}
	if e.config.IBflags&IBemojiDisabled == 0 && keyVal == IBUS_Colon && e.isEmojiLTOpened == false {
		e.resetBuffer()
		e.isEmojiLTOpened = true
		e.openEmojiList()
		return true, nil
	}
	if e.isInputModeLTOpened {
		return e.ltProcessKeyEvent(keyVal, keyCode, state)
	}
	if e.isEmojiLTOpened {
		return e.emojiProcessKeyEvent(keyVal, keyCode, state)
	}
	if e.inPreeditList() {
		return e.preeditProcessKeyEvent(keyVal, keyCode, state)
	}
	if e.inBackspaceWhiteList() {
		return e.bsProcessKeyEvent(keyVal, keyCode, state)
	}
	return e.preeditProcessKeyEvent(keyVal, keyCode, state)
}

func (e *IBusBambooEngine) FocusIn() *dbus.Error {
	log.Print("FocusIn.")
	var oldWmClasses = e.wmClasses
	e.wmClasses = x11GetFocusWindowClass()
	fmt.Printf("WM_CLASS=(%s)\n", e.wmClasses)

	e.RegisterProperties(e.propList)
	e.RequireSurroundingText()
	e.isFocusOut = false
	if oldWmClasses != e.wmClasses {
		e.firstTimeSendingBS = true
		e.resetBuffer()
		x11ClipboardReset()
	}

	return nil
}

func (e *IBusBambooEngine) FocusOut() *dbus.Error {
	log.Print("FocusOut.")
	e.isFocusOut = true
	//e.wmClasses = ""
	return nil
}

func (e *IBusBambooEngine) Reset() *dbus.Error {
	fmt.Print("Reset.")
	return nil
}

func (e *IBusBambooEngine) Enable() *dbus.Error {
	fmt.Print("Enable.")
	e.RequireSurroundingText()
	return nil
}

//@method(in_signature="vuu")
func (e *IBusBambooEngine) SetSurroundingText(text dbus.Variant, cursorPos uint32, anchorPos uint32) *dbus.Error {
	if e.getRawKeyLen() > 0 {
		return nil
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if !e.inPreeditList() && e.inBackspaceWhiteList() {
		var str = reflect.ValueOf(reflect.ValueOf(text.Value()).Index(2).Interface()).String()
		var s = []rune(str)
		if len(s) < int(cursorPos) {
			return nil
		}
		fmt.Println("Surrounding Text: ", string(s[:cursorPos]))
		e.preeditor.Reset()
		e.preeditor.ProcessString(string(s[:cursorPos]), bamboo.EnglishMode)
	}
	return nil
}

func (e *IBusBambooEngine) Disable() *dbus.Error {
	fmt.Print("Disable.")
	x11ClipboardExit()
	stopMouseTracking()
	return nil
}

func (e *IBusBambooEngine) PageUp() *dbus.Error {
	if e.isEmojiLTOpened && e.emojiLookupTable.PageUp() {
		e.updateEmojiLookupTable()
	}
	if e.isInputModeLTOpened && e.inputModeLookupTable.PageUp() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) PageDown() *dbus.Error {
	if e.isEmojiLTOpened && e.emojiLookupTable.PageDown() {
		e.updateEmojiLookupTable()
	}
	if e.isInputModeLTOpened && e.inputModeLookupTable.PageDown() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) CursorUp() *dbus.Error {
	if e.isEmojiLTOpened && e.emojiLookupTable.CursorUp() {
		e.updateEmojiLookupTable()
	}
	if e.isInputModeLTOpened && e.inputModeLookupTable.CursorUp() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) CursorDown() *dbus.Error {
	if e.isEmojiLTOpened && e.emojiLookupTable.CursorDown() {
		e.updateEmojiLookupTable()
	}
	if e.isInputModeLTOpened && e.inputModeLookupTable.CursorDown() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) CandidateClicked(index uint32, button uint32, state uint32) *dbus.Error {
	if e.isEmojiLTOpened && e.emojiLookupTable.SetCursorPosInCurrentPage(index) {
		e.commitEmojiCandidate()
		e.closeEmojiCandidates()
	}
	if e.isInputModeLTOpened && e.inputModeLookupTable.SetCursorPos(index) {
		e.commitInputModeCandidate()
		e.closeInputModeCandidates()
	}
	return nil
}

func (e *IBusBambooEngine) SetCapabilities(cap uint32) *dbus.Error {
	e.capabilities = cap
	return nil
}

func (e *IBusBambooEngine) SetCursorLocation(x int32, y int32, w int32, h int32) *dbus.Error {
	e.zeroLocation = x == 0 && y == 0 && w == 0 && h == 0
	return nil
}

func (e *IBusBambooEngine) SetContentType(purpose uint32, hints uint32) *dbus.Error {
	return nil
}

//@method(in_signature="su")
func (e *IBusBambooEngine) PropertyActivate(propName string, propState uint32) *dbus.Error {
	if propName == PropKeyAbout {
		exec.Command("xdg-open", HomePage).Start()
		return nil
	}
	if propName == PropKeyVnConvert {
		exec.Command("xdg-open", VnConvertPage).Start()
		return nil
	}
	if propName == PropKeyBambooConfiguration {
		exec.Command("xdg-open", getConfigPath(e.engineName)).Start()
		return nil
	}
	if propName == PropKeyMacroTable {
		OpenMactabFile(e.engineName)
		return nil
	}

	turnSpellChecking := func(on bool) {
		if on {
			e.config.IBflags |= IBspellChecking
			e.config.IBflags |= IBautoNonVnRestore
			if e.config.IBflags&IBspellCheckingWithDicts == 0 {
				e.config.IBflags |= IBspellCheckingWithRules
			}
		} else {
			e.config.IBflags &= ^IBspellChecking
			e.config.IBflags &= ^IBautoNonVnRestore
			e.config.IBflags &= ^IBautoCommitWithVnNotMatch
			e.config.IBflags &= ^IBautoCommitWithVnFullMatch
		}
	}

	if propName == PropKeyEmojiEnabled {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags &= ^IBemojiDisabled
		} else {
			e.config.IBflags |= IBemojiDisabled
		}
	}

	if propName == PropKeyStdToneStyle {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.Flags |= bamboo.EstdToneStyle
		} else {
			e.config.Flags &= ^bamboo.EstdToneStyle
		}
	}
	if propName == PropKeyFreeToneMarking {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.Flags |= bamboo.EfreeToneMarking
		} else {
			e.config.Flags &= ^bamboo.EfreeToneMarking
		}
	}
	if propName == PropKeySpellingChecking {
		if propState == ibus.PROP_STATE_CHECKED {
			turnSpellChecking(true)
		} else {
			turnSpellChecking(false)
		}
	}
	if propName == PropKeySpellCheckingByRules {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBspellCheckingWithRules
			turnSpellChecking(true)
		} else {
			e.config.IBflags &= ^IBspellCheckingWithRules
		}
	}
	if propName == PropKeySpellCheckingByDicts {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBspellCheckingWithDicts
			turnSpellChecking(true)
		} else {
			e.config.IBflags &= ^IBspellCheckingWithDicts
		}
	}
	if propName == PropKeyAutoCommitWithVnNotMatch {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBautoCommitWithVnNotMatch
		} else {
			e.config.IBflags &= ^IBautoCommitWithVnNotMatch
		}
	}
	if propName == PropKeyAutoCommitWithVnFullMatch {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBautoCommitWithVnFullMatch
		} else {
			e.config.IBflags &= ^IBautoCommitWithVnFullMatch
		}
	}
	if propName == PropKeyAutoCommitWithVnWordBreak {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBautoCommitWithVnWordBreak
		} else {
			e.config.IBflags &= ^IBautoCommitWithVnWordBreak
		}
	}
	if propName == PropKeyAutoCommitWithMouseMovement {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBautoCommitWithMouseMovement
			startMouseTracking()
		} else {
			e.config.IBflags &= ^IBautoCommitWithMouseMovement
			stopMouseTracking()
		}
	}
	if propName == PropKeyAutoCommitWithDelay {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBautoCommitWithDelay
		} else {
			e.config.IBflags &= ^IBautoCommitWithDelay
		}
	}
	if propName == PropKeyMacroEnabled {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBmarcoEnabled
			e.config.IBflags &= ^IBautoCommitWithVnNotMatch
			e.config.IBflags &= ^IBautoCommitWithVnFullMatch
			e.config.IBflags &= ^IBautoCommitWithVnWordBreak
			e.macroTable.Enable(e.engineName)
		} else {
			e.config.IBflags &= ^IBmarcoEnabled
			e.macroTable.Disable()
		}
	}
	if propName == PropKeyInvisibilityPreedit {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBpreeditInvisibility
		} else {
			e.config.IBflags &= ^IBpreeditInvisibility
		}
	}
	if propName == PropKeyFakeBackspace {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBfakeBackspaceEnabled
		} else {
			e.config.IBflags &= ^IBfakeBackspaceEnabled
		}
	}

	if propName == PropKeyDisableInputLookupTable {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags &= ^IBinputLookupTableDisabled
		} else {
			e.config.IBflags |= IBinputLookupTableDisabled
		}
	}

	var charset, foundCs = getCharsetFromPropKey(propName)
	if foundCs && isValidCharset(charset) && propState == ibus.PROP_STATE_CHECKED {
		e.config.OutputCharset = charset
	}
	if _, found := e.config.InputMethodDefinitions[propName]; found && propState == ibus.PROP_STATE_CHECKED {
		e.config.InputMethod = propName
	}
	SaveConfig(e.config, e.engineName)
	e.propList = GetPropListByConfig(e.config)

	var inputMethod = bamboo.ParseInputMethod(e.config.InputMethodDefinitions, e.config.InputMethod)
	e.preeditor = bamboo.NewEngine(inputMethod, e.config.Flags)
	e.RegisterProperties(e.propList)
	return nil
}
