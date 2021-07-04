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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

var dictionary = map[string]bool{}
var emojiTrie = NewTrie()

func GetIBusEngineCreator() func(*dbus.Conn, string) dbus.ObjectPath {
	go keyPressCapturing()

	return func(conn *dbus.Conn, ngName string) dbus.ObjectPath {
		var engineName = strings.ToLower(ngName)
		var engine = new(IBusBambooEngine)
		var config = loadConfig(engineName)
		var objectPath = dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/%s/%d", engineName, time.Now().UnixNano()))
		var inputMethod = bamboo.ParseInputMethod(config.InputMethodDefinitions, config.InputMethod)
		engine.Engine = ibus.BaseEngine(conn, objectPath)
		engine.engineName = engineName
		engine.preeditor = bamboo.NewEngine(inputMethod, config.Flags)
		engine.config = loadConfig(engineName)
		engine.propList = GetPropListByConfig(config)
		ibus.PublishEngine(conn, objectPath, engine)
		go engine.init()

		return objectPath
	}
}

const KeypressDelayMs = 10

func (e *IBusBambooEngine) init() {
	e.emoji = NewEmojiEngine()
	if e.macroTable == nil {
		e.macroTable = NewMacroTable()
		if e.config.IBflags&IBmacroEnabled != 0 {
			e.macroTable.Enable(e.engineName)
		}
	}
	if e.config.IBflags&IBspellCheckWithDicts != 0 && len(dictionary) == 0 {
		dictionary, _ = loadDictionary(DictVietnameseCm)
	}
	if e.config.IBflags&IBemojiDisabled == 0 && emojiTrie != nil && len(emojiTrie.Children) == 0 {
		emojiTrie, _ = loadEmojiOne(DictEmojiOne)
	}
	keyPressHandler = e.keyPressHandler

	if e.config.IBflags&IBmouseCapturing != 0 {
		startMouseCapturing()
	}
	startMouseRecording()
	var mouseMutex sync.Mutex
	onMouseMove = func() {
		mouseMutex.Lock()
		defer mouseMutex.Unlock()
		if e.checkInputMode(preeditIM) {
			if e.getRawKeyLen() == 0 {
				return
			}
			e.commitPreedit(e.getPreeditString())
		}
	}
	onMouseClick = func() {
		mouseMutex.Lock()
		defer mouseMutex.Unlock()
		if e.isEmojiLTOpened {
			e.refreshEmojiCandidate()
		} else {
			e.resetFakeBackspace()
			e.resetBuffer()
			e.keyPressDelay = KeypressDelayMs
			if e.capabilities&IBusCapSurroundingText != 0 {
				//e.ForwardKeyEvent(IBUS_Shift_R, XK_Shift_R-8, 0)
				x11SendShiftR()
				e.isSurroundingTextReady = true
				e.keyPressDelay = KeypressDelayMs * 10
			}
		}
	}
}

var keyPressHandler = func(keyVal, keyCode, state uint32) {}
var keyPressChan = make(chan [3]uint32, 100)
var isProcessing bool

func keyPressCapturing() {
	for keyEvents := range keyPressChan {
		isProcessing = true
		var keyVal, keyCode, state = keyEvents[0], keyEvents[1], keyEvents[2]
		keyPressHandler(keyVal, keyCode, state)
		isProcessing = false
	}
}

func (e *IBusBambooEngine) resetBuffer() {
	if e.getRawKeyLen() == 0 {
		return
	}
	if e.checkInputMode(preeditIM) {
		e.commitPreedit(e.getPreeditString())
	} else {
		e.preeditor.Reset()
	}
}

func (e *IBusBambooEngine) checkWmClass(newId string) {
	if e.wmClasses != newId {
		e.wmClasses = newId
		e.resetBuffer()
		e.resetFakeBackspace()
	}
}

func (e *IBusBambooEngine) processShiftKey(keyVal, state uint32) bool {
	if keyVal == IBusShiftL || keyVal == IBusShiftR {
		// when press one Shift key
		if state&IBusShiftMask != 0 && state&IBusReleaseMask != 0 &&
			e.config.IBflags&IBimQuickSwitchEnabled != 0 && !e.lastKeyWithShift {
			e.englishMode = !e.englishMode
			notify(e.englishMode)
			e.resetBuffer()
		}
		return true
	}
	return false
}

func (e *IBusBambooEngine) toUpper(keyRune rune) rune {
	var keyMapping = map[rune]rune{
		'[': '{',
		']': '}',
		'{': '[',
		'}': ']',
	}

	if upperSpecialKey, found := keyMapping[keyRune]; found && inKeyList(e.preeditor.GetInputMethod().AppendingKeys, keyRune) {
		keyRune = upperSpecialKey
	}
	return keyRune
}

func (e *IBusBambooEngine) updateLastKeyWithShift(keyVal, state uint32) {
	if e.canProcessKey(keyVal) {
		e.lastKeyWithShift = state&IBusShiftMask != 0
	} else {
		e.lastKeyWithShift = false
	}
}

func (e *IBusBambooEngine) isIgnoredKey(keyVal, keyCode, state uint32) bool {
	if state&IBusReleaseMask != 0 {
		//Ignore key-up event
		return true
	}
	if keyVal == IBusCapsLock {
		return true
	}
	if e.checkInputMode(usModeIM) {
		return true
	}
	if e.checkInputMode(usIM) {
		if e.isInputModeLTOpened || e.isInputModeShortcutKeyPressed(keyVal, keyCode, state) {
			return false
		}
		return true
	}
	return false
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.preeditor.GetProcessedString(bamboo.EnglishMode | bamboo.FullText))
}

func (e *IBusBambooEngine) getInputMode() int {
	if e.getWmClass() != "" {
		if im, ok := e.config.InputModeMapping[e.getWmClass()]; ok && imLookupTable[im] != "" {
			return im
		}
	}
	if imLookupTable[e.config.DefaultInputMode] != "" {
		return e.config.DefaultInputMode
	}
	return preeditIM
}

func (e *IBusBambooEngine) openLookupTable() {
	var wmClasses = strings.Split(e.getWmClass(), ":")
	var wmClass = e.getWmClass()
	if len(wmClasses) == 2 {
		wmClass = wmClasses[1]
	}

	e.UpdateAuxiliaryText(ibus.NewText("Nhấn (1/2/3/4/5/6/7/8) để lưu tùy chọn của bạn"), true)

	lt := ibus.NewLookupTable()
	lt.PageSize = uint32(len(imLookupTable))
	lt.Orientation = IBusOrientationVertical
	for im := 1; im <= len(imLookupTable); im++ {
		if e.getInputMode() == im {
			lt.AppendLabel("*")
			lt.SetCursorPos(uint32(im - 1))
		} else {
			lt.AppendLabel(strconv.Itoa(im))
		}
		if im == usIM || im == usModeIM {
			lt.AppendCandidate(imLookupTable[im] + " (" + wmClass + ")")
		} else {
			lt.AppendCandidate(imLookupTable[im])
		}
	}
	e.inputModeLookupTable = lt
	e.UpdateLookupTable(lt, true)
}

func (e *IBusBambooEngine) ltProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var wmClasses = e.getWmClass()
	//e.HideLookupTable()
	fmt.Printf("keyCode 0x%04x keyval 0x%04x | %c\n", keyCode, keyVal, rune(keyVal))
	//e.HideAuxiliaryText()
	if wmClasses == "" {
		return true, nil
	}
	if e.isInputModeShortcutKeyPressed(keyVal, keyCode, state) {
		e.closeInputModeCandidates()
		return false, nil
	}
	var keyRune = rune(keyVal)
	if keyVal == IBusLeft || keyVal == IBusUp {
		e.CursorUp()
		return true, nil
	} else if keyVal == IBusRight || keyVal == IBusDown {
		e.CursorDown()
		return true, nil
	} else if keyVal == IBusPageUp {
		e.PageUp()
		return true, nil
	} else if keyVal == IBusPageDown {
		e.PageDown()
		return true, nil
	}
	if keyVal == IBusReturn {
		e.commitInputModeCandidate()
		e.closeInputModeCandidates()
		return true, nil
	}
	if keyRune >= '1' && keyRune <= '7' {
		if pos, err := strconv.Atoi(string(keyRune)); err == nil {
			if e.inputModeLookupTable.SetCursorPos(uint32(pos - 1)) {
				e.commitInputModeCandidate()
				e.closeInputModeCandidates()
				return true, nil
			} else {
				e.closeInputModeCandidates()
			}
		}
	}
	e.closeInputModeCandidates()
	return false, nil
}

func (e *IBusBambooEngine) commitInputModeCandidate() {
	var im = e.inputModeLookupTable.CursorPos + 1
	e.config.InputModeMapping[e.getWmClass()] = int(im)

	saveConfig(e.config, e.engineName)
	e.propList = GetPropListByConfig(e.config)
	e.RegisterProperties(e.propList)
}

func (e *IBusBambooEngine) closeInputModeCandidates() {
	e.inputModeLookupTable = nil
	e.UpdateLookupTable(ibus.NewLookupTable(), true) // workaround for issue #18
	e.HidePreeditText()
	e.HideLookupTable()
	e.HideAuxiliaryText()
	e.isInputModeLTOpened = false
}

func (e *IBusBambooEngine) updateInputModeLT() {
	var visible = len(e.inputModeLookupTable.Candidates) > 0
	e.UpdateLookupTable(e.inputModeLookupTable, visible)
}

func (e *IBusBambooEngine) isValidState(state uint32) bool {
	if state&IBusControlMask != 0 ||
		state&IBusMod1Mask != 0 ||
		state&IBusIgnoredMask != 0 ||
		state&IBusSuperMask != 0 ||
		state&IBusHyperMask != 0 ||
		state&IBusMetaMask != 0 {
		return false
	}
	return true
}

func (e *IBusBambooEngine) getCommitText(keyVal, keyCode, state uint32) (string, bool) {
	var keyRune = rune(keyVal)
	oldText := e.getPreeditString()
	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBusLockMask != 0 {
			keyRune = e.toUpper(keyRune)
		}
		e.preeditor.ProcessKey(keyRune, e.getBambooInputMode())
		if inKeyList(e.preeditor.GetInputMethod().AppendingKeys, keyRune) {
			var newText string
			if e.shouldFallbackToEnglish(true) {
				newText = e.getProcessedString(bamboo.EnglishMode)
			} else {
				newText = e.getProcessedString(bamboo.VietnameseMode)
			}
			if fullSeq := e.preeditor.GetProcessedString(bamboo.VietnameseMode); len(fullSeq) > 0 && rune(fullSeq[len(fullSeq)-1]) == keyRune {
				// [[ => [
				var ret = e.getPreeditString()
				e.preeditor.Reset()
				return ret, true
			} else if newText != "" && keyRune == rune(newText[len(newText)-1]) {
				// f] => f]
				e.preeditor.Reset()
				return oldText + string(keyRune), true
			} else {
				// ] => o?
				return e.getPreeditString(), false
			}
		} else if e.config.IBflags&IBmacroEnabled != 0 {
			return e.getProcessedString(bamboo.PunctuationMode), false
		} else {
			return e.getPreeditString(), false
		}
	} else if bamboo.IsWordBreakSymbol(keyRune) {
		// restore key strokes by pressing Shift + Space
		if keyVal == IBusSpace && state&IBusShiftMask != 0 &&
			e.config.IBflags&IBrestoreKeyStrokesEnabled != 0 && !e.lastKeyWithShift {
			if bamboo.HasAnyVietnameseRune(oldText) {
				commitText := e.preeditor.GetProcessedString(bamboo.EnglishMode)
				e.preeditor.RestoreLastWord()
				return commitText, false
			}
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return oldText + string(keyRune), true
		}
		// macro processing
		if e.config.IBflags&IBmacroEnabled != 0 {
			var keyS = string(keyRune)
			if keyVal == IBusSpace && e.macroTable.HasKey(oldText) {
				e.preeditor.Reset()
				return e.expandMacro(oldText) + keyS, keyVal == IBusSpace
			} else {
				e.preeditor.ProcessKey(keyRune, e.getBambooInputMode())
				return oldText + keyS, keyVal == IBusSpace
			}
		}
		if bamboo.HasAnyVietnameseRune(oldText) && e.mustFallbackToEnglish() {
			e.preeditor.RestoreLastWord()
			newText := e.preeditor.GetProcessedString(bamboo.EnglishMode) + string(keyRune)
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return newText, true
		}
		e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		return oldText + string(keyRune), true
	}
	return "", true
}

func (e *IBusBambooEngine) commitMacroText(keyRune rune) bool {
	if e.config.IBflags&IBmacroEnabled == 0 {
		return false
	}
	var keyS = string(keyRune)
	var text = e.preeditor.GetProcessedString(bamboo.PunctuationMode)
	if e.macroTable.HasKey(text) {
		e.commitPreedit(e.expandMacro(text) + keyS)
		return true
	} else if e.macroTable.HasKey(text + keyS) {
		e.preeditor.ProcessKey(keyRune, e.getBambooInputMode())
		e.updatePreedit(text + keyS)
		return true
	}
	return false
}

func (e *IBusBambooEngine) getMacroText() (bool, string) {
	if e.config.IBflags&IBmacroEnabled == 0 {
		return false, ""
	}
	var text = e.preeditor.GetProcessedString(bamboo.PunctuationMode)
	if e.macroTable.HasKey(text) {
		return true, e.expandMacro(text)
	}
	return false, ""
}

func (e *IBusBambooEngine) getFakeBackspace() int {
	return e.nFakeBackSpace
}

var mtx sync.Mutex

func (e *IBusBambooEngine) setFakeBackspace(n int) {
	mtx.Lock()
	e.nFakeBackSpace = n
	mtx.Unlock()
}

func (e *IBusBambooEngine) addFakeBackspace(n int) {
	mtx.Lock()
	e.nFakeBackSpace += n
	mtx.Unlock()
}

func (e *IBusBambooEngine) canProcessKey(keyVal uint32) bool {
	var keyRune = rune(keyVal)
	if keyVal == IBusSpace || keyVal == IBusBackSpace || bamboo.IsWordBreakSymbol(keyRune) {
		return true
	}
	if ok, _ := e.getMacroText(); ok && keyVal == IBusTab {
		return true
	}
	return e.preeditor.CanProcessKey(keyRune)
}

func (e *IBusBambooEngine) inBackspaceWhiteList() bool {
	var inputMode = e.getInputMode()
	for _, im := range imBackspaceList {
		if im == inputMode {
			return true
		}
	}
	return false
}

func (e *IBusBambooEngine) inBrowserList() bool {
	return inStringList(DefaultBrowserList, e.getWmClass())
}

func (e *IBusBambooEngine) getWmClass() string {
	return e.wmClasses
}

func (e *IBusBambooEngine) getLatestWmClass() string {
	var wmClass string
	if isGnome {
		wmClass, _ = gnomeGetFocusWindowClass()
	} else if isWayland {
		wmClass = wlAppId
	}
	if wmClass == "" {
		wmClass = x11GetFocusWindowClass()
	}
	return wmClass
}

func (e *IBusBambooEngine) checkInputMode(im int) bool {
	return e.getInputMode() == im
}

func notify(enMode bool) {
	var title = "Vietnamese"
	var msg = "Press Shift to switch to English"
	if enMode {
		title = "English"
		msg = "Press Shift to switch to Vietnamese"
	}
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println(err)
		return
	}
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(281025),
		"", title, msg, []string{}, map[string]dbus.Variant{}, int32(3000))
	if call.Err != nil {
		fmt.Println(call.Err)
	}
}
