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
	"log"
	"os/exec"
	"reflect"
	"sync"

	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

type IBusBambooEngine struct {
	sync.Mutex
	ibus.Engine
	preeditor              bamboo.IEngine
	engineName             string
	config                 *Config
	propList               *ibus.PropList
	englishMode            bool
	macroTable             *MacroTable
	wmClasses              string
	isInputModeLTOpened    bool
	isEmojiLTOpened        bool
	emojiLookupTable       *ibus.LookupTable
	inputModeLookupTable   *ibus.LookupTable
	capabilities           uint32
	keyPressDelay          int
	nFakeBackSpace         int
	isFirstTimeSendingBS   bool
	emoji                  *EmojiEngine
	isSurroundingTextReady bool
	lastKeyWithShift       bool
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
	if e.checkInputMode(usIM) {
		if e.isInputModeLTOpened || keyVal == IBusOpenLookupTable {
			// return false, nil
		} else {
			return false, nil
		}
	}
	if e.processShiftKey(keyVal, state) {
		return true, nil
	}
	if e.isIgnoredKey(keyVal, state) {
		return false, nil
	}
	log.Printf("ProcessKeyEvent >  %c | keyCode 0x%04x keyVal 0x%04x | %d\n", rune(keyVal), keyCode, keyVal, len(keyPressChan))
	if e.config.IBflags&IBinputModeLookupTableEnabled != 0 && keyVal == IBusOpenLookupTable && !e.isInputModeLTOpened && e.wmClasses != "" {
		e.resetBuffer()
		e.isInputModeLTOpened = true
		e.lastKeyWithShift = true
		e.openLookupTable()
		return true, nil
	}
	if e.config.IBflags&IBemojiDisabled == 0 && keyVal == IBusColon && !e.isEmojiLTOpened {
		e.resetBuffer()
		e.isEmojiLTOpened = true
		e.lastKeyWithShift = true
		e.openEmojiList()
		return true, nil
	}
	if e.isInputModeLTOpened {
		return e.ltProcessKeyEvent(keyVal, keyCode, state)
	}
	if e.isEmojiLTOpened {
		return e.emojiProcessKeyEvent(keyVal, keyCode, state)
	}
	if e.englishMode {
		e.updateLastKeyWithShift(keyVal, state)
		return false, nil
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
	if oldWmClasses != e.wmClasses {
		e.resetBuffer()
		e.resetFakeBackspace()
	}
	return nil
}

func (e *IBusBambooEngine) FocusOut() *dbus.Error {
	log.Print("FocusOut.")
	//e.wmClasses = ""
	return nil
}

func (e *IBusBambooEngine) Reset() *dbus.Error {
	fmt.Print("Reset.\n")
	if e.checkInputMode(preeditIM) {
		e.commitPreedit(e.getPreeditString())
	}
	return nil
}

func (e *IBusBambooEngine) Enable() *dbus.Error {
	fmt.Print("Enable.")
	e.RequireSurroundingText()
	return nil
}

func (e *IBusBambooEngine) Disable() *dbus.Error {
	fmt.Print("Disable.")
	return nil
}

//@method(in_signature="vuu")
func (e *IBusBambooEngine) SetSurroundingText(text dbus.Variant, cursorPos uint32, anchorPos uint32) *dbus.Error {
	if !e.isSurroundingTextReady {
		//fmt.Println("Surrounding Text is not ready yet.")
		return nil
	}
	e.Lock()
	defer func() {
		e.Unlock()
		e.isSurroundingTextReady = false
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if e.inBackspaceWhiteList() {
		var str = reflect.ValueOf(reflect.ValueOf(text.Value()).Index(2).Interface()).String()
		var s = []rune(str)
		if len(s) < int(cursorPos) {
			return nil
		}
		var cs = s[:cursorPos]
		fmt.Println("Surrounding Text: ", string(cs))
		e.preeditor.Reset()
		for i := len(cs) - 1; i >= 0; i-- {
			// workaround for spell checking
			if bamboo.IsPunctuationMark(cs[i]) && e.preeditor.CanProcessKey(cs[i]) {
				cs[i] = ' '
			}
			e.preeditor.ProcessKey(cs[i], bamboo.EnglishMode|bamboo.InReverseOrder)
		}
	}
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
	if e.isEmojiLTOpened && e.updateCursorPosInEmojiTable(index) {
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
	if propName == PropKeyVnCharsetConvert {
		exec.Command("xdg-open", CharsetConvertPage).Start()
		return nil
	}
	if propName == PropKeyConfiguration {
		exec.Command("xdg-open", getConfigPath(e.engineName)).Start()
		return nil
	}
	if propName == PropKeyMacroTable {
		OpenMactabFile(e.engineName)
		return nil
	}

	turnSpellChecking := func(on bool) {
		if on {
			e.config.IBflags |= IBspellCheckEnabled
			e.config.IBflags |= IBautoNonVnRestore
			if e.config.IBflags&IBspellCheckWithDicts == 0 {
				e.config.IBflags |= IBspellCheckWithRules
			}
		} else {
			e.config.IBflags &= ^IBspellCheckEnabled
			e.config.IBflags &= ^IBautoNonVnRestore
		}
	}

	if propName == PropKeyEmojiEnabled {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags &= ^IBemojiDisabled
			emojiTrie, _ = loadEmojiOne(DictEmojiOne)
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
	if propName == PropKeyEnableSpellCheck {
		if propState == ibus.PROP_STATE_CHECKED {
			turnSpellChecking(true)
		} else {
			turnSpellChecking(false)
		}
	}
	if propName == PropKeySpellCheckByRules {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBspellCheckWithRules
			turnSpellChecking(true)
		} else {
			e.config.IBflags &= ^IBspellCheckWithRules
		}
	}
	if propName == PropKeySpellCheckByDicts {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBspellCheckWithDicts
			turnSpellChecking(true)
			dictionary, _ = loadDictionary(DictVietnameseCm)
		} else {
			e.config.IBflags &= ^IBspellCheckWithDicts
		}
	}
	if propName == PropKeyMouseCapturing {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBmouseCapturing
			startMouseCapturing()
		} else {
			e.config.IBflags &= ^IBmouseCapturing
			stopMouseCapturing()
		}
	}
	if propName == PropKeyMacroEnabled {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBmacroEnabled
			e.config.IBflags |= IBautoCapitalizeMacro
			e.macroTable.Enable(e.engineName)
		} else {
			e.config.IBflags &= ^IBmacroEnabled
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
	if propName == PropKeyPreeditElimination {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBpreeditElimination
		} else {
			e.config.IBflags &= ^IBpreeditElimination
		}
	}
	if propName == PropKeyRestoreKeyStrokes {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBrestoreKeyStrokesEnabled
		} else {
			e.config.IBflags &= ^IBrestoreKeyStrokesEnabled
		}
	}
	if propName == PropKeyInputModeLookupTable {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBinputModeLookupTableEnabled
		} else {
			e.config.IBflags &= ^IBinputModeLookupTableEnabled
		}
	}
	if propName == PropKeyIMQuickSwitchEnabled {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBimQuickSwitchEnabled
		} else {
			e.config.IBflags &= ^IBimQuickSwitchEnabled
		}
		e.englishMode = false
	}
	if propName == PropKeyAutoCapitalizeMacro {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBautoCapitalizeMacro
		} else {
			e.config.IBflags &= ^IBautoCapitalizeMacro
		}
	}

	var charset, foundCs = getCharsetFromPropKey(propName)
	if foundCs && isValidCharset(charset) && propState == ibus.PROP_STATE_CHECKED {
		e.config.OutputCharset = charset
	}
	if _, found := e.config.InputMethodDefinitions[propName]; found && propState == ibus.PROP_STATE_CHECKED {
		e.config.InputMethod = propName
	}
	if propName != "-" {
		saveConfig(e.config, e.engineName)
	}
	e.propList = GetPropListByConfig(e.config)

	var inputMethod = bamboo.ParseInputMethod(e.config.InputMethodDefinitions, e.config.InputMethod)
	e.preeditor = bamboo.NewEngine(inputMethod, e.config.Flags)
	e.RegisterProperties(e.propList)
	return nil
}
