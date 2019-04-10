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
	"sync"
)

type IBusBambooEngine struct {
	sync.Mutex
	ibus.Engine
	preeditor           bamboo.IEngine
	zeroLocation        bool
	engineName          string
	config              *Config
	propList            *ibus.PropList
	mode                bamboo.Mode
	ignorePreedit       bool
	macroTable          *MacroTable
	dictionary          map[string]bool
	wmClasses           string
	isLookupTableOpened bool
	isEmojiTableOpened  bool
	emojiLookupTable    *ibus.LookupTable
	capabilities        uint32
	nFakeBackSpace      int
	firstTimeSendingBS  bool
	emoji               *BambooEmoji
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
	if keyVal == IBUS_OpenLookupTable && e.isLookupTableOpened == false {
		e.resetBuffer()
		e.isLookupTableOpened = true
		e.openLookupTable()
		return true, nil
	}
	if e.config.IBflags&IBemojiDisabled == 0 && keyVal == IBUS_Colon && e.isEmojiTableOpened == false {
		e.resetBuffer()
		e.isEmojiTableOpened = true
		e.openEmojiList()
		return true, nil
	}
	if e.isLookupTableOpened {
		e.isLookupTableOpened = false
		return e.ltProcessKeyEvent(keyVal, keyCode, state)
	}
	if e.isEmojiTableOpened {
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
	e.HidePreeditText()
	e.firstTimeSendingBS = true
	if oldWmClasses != e.wmClasses {
		e.resetBuffer()
		x11ClipboardReset()
	}

	return nil
}

func (e *IBusBambooEngine) FocusOut() *dbus.Error {
	log.Print("FocusOut.")
	//e.wmClasses = ""
	return nil
}

func (e *IBusBambooEngine) Reset() *dbus.Error {
	fmt.Print("Reset.")
	return nil
}

func (e *IBusBambooEngine) Enable() *dbus.Error {
	fmt.Print("Enable.")
	if e.config.IBflags&IBautoCommitWithMouseMovement != 0 {
		mouseCaptureInit()
	}
	if e.config.IBflags&IBmarcoEnabled != 0 {
		e.macroTable.Enable(e.engineName)
	}
	return nil
}

func (e *IBusBambooEngine) Disable() *dbus.Error {
	fmt.Print("Disable.")
	x11ClipboardExit()
	mouseCaptureExit()
	return nil
}

func (e *IBusBambooEngine) PageUp() *dbus.Error {
	if e.isEmojiTableOpened && e.emojiLookupTable.PageUp() {
		e.emojiUpdateLookupTable()
	}
	return nil
}

func (e *IBusBambooEngine) PageDown() *dbus.Error {
	if e.isEmojiTableOpened && e.emojiLookupTable.PageDown() {
		e.emojiUpdateLookupTable()
	}
	return nil
}

func (e *IBusBambooEngine) CursorUp() *dbus.Error {
	if e.isEmojiTableOpened && e.emojiLookupTable.CursorUp() {
		e.emojiUpdateLookupTable()
	}
	return nil
}

func (e *IBusBambooEngine) CursorDown() *dbus.Error {
	if e.isEmojiTableOpened && e.emojiLookupTable.CursorDown() {
		e.emojiUpdateLookupTable()
	}
	return nil
}

func (e *IBusBambooEngine) CandidateClicked(index uint32, button uint32, state uint32) *dbus.Error {
	if e.isEmojiTableOpened && e.emojiLookupTable.SetCursorPosInCurrentPage(index) {
		e.commitEmojiCandidate()
		e.closeEmojiCandidates()
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
			mouseCaptureInit()
		} else {
			e.config.IBflags &= ^IBautoCommitWithMouseMovement
			mouseCaptureExit()
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
	if propName == PropKeyPreeditInvisibility {
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
	e.preeditor = bamboo.NewEngine(inputMethod, e.config.Flags, e.dictionary)
	e.RegisterProperties(e.propList)
	return nil
}
