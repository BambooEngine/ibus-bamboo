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

// / This is public API for IBus, these functionalities has to be methods to work with IBus
package main

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"strconv"

	"github.com/BambooEngine/bamboo-core"
	ibus "github.com/BambooEngine/goibus"
	"github.com/godbus/dbus/v5"

	"ibus-bamboo/config"
	"ibus-bamboo/emoji"
	"ibus-bamboo/introspectors"
	"ibus-bamboo/mactab"
	"ibus-bamboo/types"
	"ibus-bamboo/ui"
)

type IBusBambooEngine struct {
	*types.IBusBambooEngineState
}

func NewIbusBambooEngine(name string, cfg *types.Config, base types.IEngine, preeditor bamboo.IEngine) *IBusBambooEngine {
	return &IBusBambooEngine{&types.IBusBambooEngineState{
		EngineName: name,
		IEngine:    base,
		Preeditor:  preeditor,
		Config:     cfg,
	}}
}

/*
*
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
	if state&IBusReleaseMask != 0 {
		// fmt.Println("Ignore key-up event")
		return false, nil
	}
	fmt.Printf("\n")
	log.Printf(">>>>ProcessKeyEvent >  %d | state %d keyVal 0x%04x | %c <<<<\n", len(keyPressChan), state, keyVal, rune(keyVal))
	if ret, retValue := e.processShortcutKey(keyVal, keyCode, state); ret {
		return retValue, nil
	}
	if e.inBackspaceWhiteList() {
		return e.bsProcessKeyEvent(keyVal, keyCode, state)
	}
	return e.preeditProcessKeyEvent(keyVal, keyCode, state)
}

func (e *IBusBambooEngine) FocusIn() *dbus.Error {
	log.Print("FocusIn.")
	var latestWm = introspectors.GetLatestAppId(introspector)
	e.checkWmClass(latestWm)
	e.RegisterProperties(e.PropList)
	e.RequireSurroundingText()
	if e.isShortcutKeyEnable(KSEmojiDialog) && emoji.Trie != nil && len(emoji.Trie.Children) == 0 {
		var err error
		emoji.Trie, err = emoji.LoadEmojiOne(DictEmojiOne)
		if err != nil {
			panic(fmt.Sprintf("failed to load emojiTrie from %s: %s", DictEmojiOne, err))
		}
	}
	if e.Config.IBflags&config.IBspellCheckWithDicts != 0 && len(dictionary) == 0 {
		dictionary, _ = loadDictionary(DictVietnameseCm)
	}
	fmt.Printf("WM_CLASS=(%s)\n", e.getWmClass())
	return nil
}

func (e *IBusBambooEngine) FocusOut() *dbus.Error {
	log.Print("FocusOut.")
	return nil
}

func (e *IBusBambooEngine) Reset() *dbus.Error {
	fmt.Print("Reset.\n")
	if e.checkInputMode(config.PreeditIM) {
		e.Preeditor.Reset()
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

// @method(in_signature="vuu")
func (e *IBusBambooEngine) SetSurroundingText(text dbus.Variant, cursorPos uint32, anchorPos uint32) *dbus.Error {
	if !e.IsSurroundingTextReady {
		//fmt.Println("Surrounding Text is not ready yet.")
		return nil
	}
	e.Lock()
	defer func() {
		e.Unlock()
		e.IsSurroundingTextReady = false
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
		e.Preeditor.Reset()
		for i := len(cs) - 1; i >= 0; i-- {
			// workaround for spell checking
			if bamboo.IsPunctuationMark(cs[i]) && e.Preeditor.CanProcessKey(cs[i]) {
				cs[i] = ' '
			}
			e.Preeditor.ProcessKey(cs[i], bamboo.EnglishMode|bamboo.InReverseOrder)
		}
	}
	return nil
}

func (e *IBusBambooEngine) PageUp() *dbus.Error {
	if e.IsEmojiLTOpened && e.EmojiLookupTable.PageUp() {
		e.updateEmojiLookupTable()
	}
	if e.IsInputModeLTOpened && e.InputModeLookupTable.PageUp() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) PageDown() *dbus.Error {
	if e.IsEmojiLTOpened && e.EmojiLookupTable.PageDown() {
		e.updateEmojiLookupTable()
	}
	if e.IsInputModeLTOpened && e.InputModeLookupTable.PageDown() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) CursorUp() *dbus.Error {
	if e.IsEmojiLTOpened && e.EmojiLookupTable.CursorUp() {
		e.updateEmojiLookupTable()
	}
	if e.IsInputModeLTOpened && e.InputModeLookupTable.CursorUp() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) CursorDown() *dbus.Error {
	if e.IsEmojiLTOpened && e.EmojiLookupTable.CursorDown() {
		e.updateEmojiLookupTable()
	}
	if e.IsInputModeLTOpened && e.InputModeLookupTable.CursorDown() {
		e.updateInputModeLT()
	}
	return nil
}

func (e *IBusBambooEngine) CandidateClicked(index uint32, button uint32, state uint32) *dbus.Error {
	if e.IsEmojiLTOpened && e.updateCursorPosInEmojiTable(index) {
		e.commitEmojiCandidate()
		e.closeEmojiCandidates()
	}
	if e.IsInputModeLTOpened && e.InputModeLookupTable.SetCursorPos(index) {
		e.commitInputModeCandidate()
		e.closeInputModeCandidates()
	}
	return nil
}

func (e *IBusBambooEngine) SetCapabilities(cap uint32) *dbus.Error {
	e.Capabilities = cap
	return nil
}

func (e *IBusBambooEngine) SetCursorLocation(x int32, y int32, w int32, h int32) *dbus.Error {
	return nil
}

func (e *IBusBambooEngine) SetContentType(purpose uint32, hints uint32) *dbus.Error {
	return nil
}

// @method(in_signature="su")
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
		ui.OpenGUI(e.EngineName)
		e.Config = config.LoadConfig(e.EngineName)
		return nil
	}
	if propName == PropKeyInputModeLookupTableShortcut {
		ui.OpenGUI(e.EngineName)
		e.Config = config.LoadConfig(e.EngineName)
		return nil
	}
	if propName == PropKeyMacroTable {
		ui.OpenGUI(e.EngineName)
		e.Config = config.LoadConfig(e.EngineName)
		return nil
	}

	turnSpellChecking := func(on bool) {
		if on {
			e.Config.IBflags |= config.IBspellCheckEnabled
			e.Config.IBflags |= config.IBautoNonVnRestore
			if e.Config.IBflags&config.IBspellCheckWithDicts == 0 {
				e.Config.IBflags |= config.IBspellCheckWithRules
			}
		} else {
			e.Config.IBflags &= ^config.IBspellCheckEnabled
			e.Config.IBflags &= ^config.IBautoNonVnRestore
		}
	}

	if propName == PropKeyStdToneStyle {
		if propState == ibus.PROP_STATE_CHECKED {
			e.Config.Flags |= bamboo.EstdToneStyle
		} else {
			e.Config.Flags &= ^bamboo.EstdToneStyle
		}
	}
	if propName == PropKeyFreeToneMarking {
		if propState == ibus.PROP_STATE_CHECKED {
			e.Config.Flags |= bamboo.EfreeToneMarking
		} else {
			e.Config.Flags &= ^bamboo.EfreeToneMarking
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
			e.Config.IBflags |= config.IBspellCheckWithRules
			turnSpellChecking(true)
		} else {
			e.Config.IBflags &= ^config.IBspellCheckWithRules
		}
	}
	if propName == PropKeySpellCheckByDicts {
		if propState == ibus.PROP_STATE_CHECKED {
			e.Config.IBflags |= config.IBspellCheckWithDicts
			turnSpellChecking(true)
			dictionary, _ = loadDictionary(DictVietnameseCm)
		} else {
			e.Config.IBflags &= ^config.IBspellCheckWithDicts
		}
	}
	if propName == PropKeyMacroEnabled {
		if propState == ibus.PROP_STATE_CHECKED {
			e.Config.IBflags |= config.IBmacroEnabled
			mactab.Enable(e.MacroTable, e.EngineName)
		} else {
			e.Config.IBflags &= ^config.IBmacroEnabled
			mactab.Disable(e.MacroTable)
		}
	}
	if propName == PropKeyPreeditInvisibility {
		if propState == ibus.PROP_STATE_CHECKED {
			e.Config.IBflags |= config.IBnoUnderline
		} else {
			e.Config.IBflags &= ^config.IBnoUnderline
		}
	}
	if propName == PropKeyPreeditElimination {
		if propState == ibus.PROP_STATE_CHECKED {
			e.Config.IBflags |= config.IBpreeditElimination
		} else {
			e.Config.IBflags &= ^config.IBpreeditElimination
		}
	}
	if propName == PropKeyAutoCapitalizeMacro {
		if propState == ibus.PROP_STATE_CHECKED {
			e.Config.IBflags |= config.IBautoCapitalizeMacro
		} else {
			e.Config.IBflags &= ^config.IBautoCapitalizeMacro
		}
		if e.Config.IBflags&config.IBmacroEnabled != 0 {
			mactab.Reload(e.MacroTable, e.EngineName, e.Config.IBflags&config.IBautoCapitalizeMacro != 0)
		}
	}

	var im, foundIm = getValueFromPropKey(propName, "InputMode")
	if foundIm && propState == ibus.PROP_STATE_CHECKED {
		e.Config.DefaultInputMode, _ = strconv.Atoi(im)
	}
	var charset, foundCs = getValueFromPropKey(propName, "OutputCharset")
	if foundCs && isValidCharset(charset) && propState == ibus.PROP_STATE_CHECKED {
		e.Config.OutputCharset = charset
	}
	if _, found := e.Config.InputMethodDefinitions[propName]; found && propState == ibus.PROP_STATE_CHECKED {
		e.Config.InputMethod = propName
	}
	if propName != "-" {
		config.SaveConfig(e.Config, e.EngineName)
	}
	e.PropList = GetPropListByConfig(e.Config)

	var inputMethod = bamboo.ParseInputMethod(e.Config.InputMethodDefinitions, e.Config.InputMethod)
	e.Preeditor = bamboo.NewEngine(inputMethod, e.Config.Flags)
	e.RegisterProperties(e.PropList)
	return nil
}
