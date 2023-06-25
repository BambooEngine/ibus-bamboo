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
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

var dictionary = map[string]bool{}
var emojiTrie = NewTrie()

func GetIBusEngineCreator() func(*dbus.Conn, string) dbus.ObjectPath {
	go keyPressCapturing()

	return func(conn *dbus.Conn, ngName string) dbus.ObjectPath {
		var ngGroupName = strings.Split(ngName, "::")[0]
		var engineName = strings.ToLower(ngGroupName)
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
		if *gui {
			engine.openShortcutsGUI()
			saveConfig(engine.config, engine.engineName)
			os.Exit(0)
		}
		go engine.init()

		return objectPath
	}
}

const KeypressDelayMs = 10

func (e *IBusBambooEngine) isShortcutKeyEnable(ski uint) bool {
	if int(ski+2) > len(e.config.Shortcuts) {
		return false
	}
	l := e.config.Shortcuts[ski : ski+2]
	return l[1] > 0
}

func (e *IBusBambooEngine) init() {
	e.emoji = NewEmojiEngine()
	if e.macroTable == nil {
		e.macroTable = NewMacroTable(e.config.IBflags&IBautoCapitalizeMacro != 0)
		if e.config.IBflags&IBmacroEnabled != 0 {
			e.macroTable.Enable(e.engineName)
		}
	}
	keyPressHandler = e.keyPressHandler

	if e.config.IBflags&IBmouseCapturing != 0 {
		startMouseCapturing()
		startMouseRecording()
	}
	var mouseMutex sync.Mutex
	onMouseMove = func() {
		mouseMutex.Lock()
		defer mouseMutex.Unlock()
		if e.checkInputMode(preeditIM) {
			if e.getRawKeyLen() == 0 {
				return
			}
			e.commitPreeditAndReset(e.getPreeditString())
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
		e.commitPreeditAndReset(e.getPreeditString())
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

func (e *IBusBambooEngine) isShortcutKeyPressed(keyVal, state uint32, shortcut uint) bool {
	if !e.isShortcutKeyEnable(shortcut) {
		return false
	}
	realState := state & IBusDefaultModMask
	lowerKey := uint32(unicode.ToLower(rune(keyVal)))
	shortcuts := e.config.Shortcuts[shortcut : shortcut+2]
	ret := shortcuts[0] == realState && shortcuts[1] == lowerKey
	// fmt.Println("...isShortcutKeyPressed=", ret, ret && !e.lastKeyWithShift, shortcuts)
	if realState == 1 && shortcut == KSViEnSwitch {
		return ret && !e.lastKeyWithShift
	}
	return ret
}

func (e *IBusBambooEngine) processShortcutKey(keyVal, keyCode, state uint32) (bool, bool) {
	if keyVal == IBusCapsLock {
		return true, false
	}
	// fmt.Println("===Process shortcut for emoji selector")
	if e.isShortcutKeyPressed(keyVal, state, KSEmojiDialog) &&
		!e.isEmojiLTOpened {
		e.resetBuffer()
		e.isEmojiLTOpened = true
		e.lastKeyWithShift = true
		e.openEmojiList()
		return true, true
	}
	if e.isEmojiLTOpened {
		return true, e.emojiProcessKeyEvent(keyVal, keyCode, state)
	}
	// fmt.Println("====== Process hexadecimal key pressed")
	if e.isShortcutKeyPressed(keyVal, state, KSHexadecimal) {
		e.resetBuffer()
		e.isInHexadecimal = true
		e.setupHexadecimalProcessKeyEvent()
		return true, true
	}
	if e.isInHexadecimal {
		if e.isShortcutKeyPressed(keyVal, state, KSHexadecimal) {
			e.closeHexadecimalInput()
			e.updateLastKeyWithShift(keyVal, state)
			return true, false
		}
		return true, e.hexadecimalProcessKeyEvent(keyVal, keyCode, state)
	}

	if e.config.DefaultInputMode == usIM {
		return true, false
	}
	// fmt.Println("===== Process restoring key strokes")
	if e.isShortcutKeyPressed(keyVal, state, KSRestoreKeyStrokes) {
		e.shouldRestoreKeyStrokes = true
		return false, false
	}
	// fmt.Println("===Process shortcut for input method switcher")
	if e.isShortcutKeyPressed(keyVal, state, KSViEnSwitch) {
		e.englishMode = !e.englishMode
		notify(e.englishMode)
		e.resetBuffer()
		return true, true
	}
	// fmt.Println("====== Process shortcut for input mode switch")
	if e.isInputModeLTOpened {
		return e.ltProcessKeyEvent(keyVal, keyCode, state)
	} else if e.isShortcutKeyPressed(keyVal, state, KSInputModeSwitch) &&
		e.getWmClass() != "" {
		e.resetBuffer()
		e.isInputModeLTOpened = true
		e.lastKeyWithShift = true
		e.openLookupTable()
		return true, true
	}

	if keyVal == IBusShiftL || keyVal == IBusShiftR {
		return true, false
	}
	if e.checkInputMode(usIM) {
		if e.isInputModeLTOpened && e.isShortcutKeyPressed(keyVal, state, KSInputModeSwitch) {
			return false, false
		}
		return true, false
	}
	if e.englishMode {
		e.updateLastKeyWithShift(keyVal, state)
		return true, false
	}
	return false, false
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
	if e.preeditor.CanProcessKey(rune(keyVal)) {
		e.lastKeyWithShift = state&IBusShiftMask != 0
	} else {
		e.lastKeyWithShift = false
	}
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.preeditor.GetProcessedString(bamboo.EnglishMode | bamboo.FullText))
}

func (e *IBusBambooEngine) runeCount() int {
	return utf8.RuneCountInString(e.getPreeditString())
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

	e.UpdateAuxiliaryText(ibus.NewText("Nhấn (1/2/3/4/5/6/7) để lưu tùy chọn của bạn"), true)

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
		if im == usIM {
			lt.AppendCandidate(imLookupTable[im] + " (" + wmClass + ")")
		} else {
			lt.AppendCandidate(imLookupTable[im])
		}
	}
	e.inputModeLookupTable = lt
	e.UpdateLookupTable(lt, true)
}

func (e *IBusBambooEngine) ltProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, bool) {
	var wmClasses = e.getWmClass()
	// e.HideLookupTable()
	// e.HideAuxiliaryText()
	if wmClasses == "" {
		return true, true
	}
	if e.isShortcutKeyPressed(keyVal, state, KSInputModeSwitch) {
		e.closeInputModeCandidates()
		return true, false
	}
	var keyRune = rune(keyVal)
	if keyVal == IBusLeft || keyVal == IBusUp {
		e.CursorUp()
		return true, true
	} else if keyVal == IBusRight || keyVal == IBusDown {
		e.CursorDown()
		return true, true
	} else if keyVal == IBusPageUp {
		e.PageUp()
		return true, true
	} else if keyVal == IBusPageDown {
		e.PageDown()
		return true, true
	}
	if keyVal == IBusReturn {
		e.commitInputModeCandidate()
		e.closeInputModeCandidates()
		return true, true
	}
	if keyRune >= '1' && keyRune <= '7' {
		if pos, err := strconv.Atoi(string(keyRune)); err == nil {
			if e.inputModeLookupTable.SetCursorPos(uint32(pos - 1)) {
				e.commitInputModeCandidate()
				e.closeInputModeCandidates()
				return true, true
			} else {
				e.closeInputModeCandidates()
			}
		}
	}
	if keyVal == IBusEscape {
		e.closeInputModeCandidates()
	}
	return true, false
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
	// restore key strokes by pressing Shift + Space
	if e.shouldRestoreKeyStrokes {
		e.shouldRestoreKeyStrokes = false
		e.preeditor.RestoreLastWord(!bamboo.HasAnyVietnameseRune(oldText))
		return e.getPreeditString(), false
	}
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
				var lastRune = rune(ret[len(ret)-1])
				var isWordBreakRune = bamboo.IsWordBreakSymbol(lastRune)
				// TODO: THIS IS HACKING
				if isWordBreakRune {
					e.preeditor.RemoveLastChar(false)
					e.preeditor.ProcessKey(' ', bamboo.EnglishMode)
				}
				return ret, isWordBreakRune
			} else if l := []rune(newText); len(l) > 0 && keyRune == l[len(l)-1] {
				// f] => f]
				var isWordBreakRune = bamboo.IsWordBreakSymbol(keyRune)
				if isWordBreakRune {
					e.preeditor.RemoveLastChar(false)
					e.preeditor.ProcessKey(' ', bamboo.EnglishMode)
				}
				return oldText + string(keyRune), isWordBreakRune
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
		var isWordBreakRune = true
		// macro processing
		if e.config.IBflags&IBmacroEnabled != 0 {
			isWordBreakRune = keyVal == IBusSpace
			var keyS = string(keyRune)
			if keyVal == IBusSpace && e.macroTable.HasKey(oldText) {
				e.preeditor.Reset()
				return e.expandMacro(oldText) + keyS, true
			}
		}
		if bamboo.HasAnyVietnameseRune(oldText) && e.mustFallbackToEnglish() {
			e.preeditor.RestoreLastWord(false)
			newText := e.preeditor.GetProcessedString(bamboo.PunctuationMode|bamboo.EnglishMode) + string(keyRune)
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return newText, isWordBreakRune
		}
		e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		return oldText + string(keyRune), isWordBreakRune
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
		e.commitPreeditAndReset(e.expandMacro(text) + keyS)
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
	wmClass = strings.Replace(wmClass, "\"", "", -1)
	return wmClass
}

func (e *IBusBambooEngine) checkInputMode(im int) bool {
	return e.getInputMode() == im
}

func notify(enMode bool) {
	var title = "Vietnamese"
	var msg = "Press Shortcut keys to switch input language"
	if enMode {
		title = "English"
		msg = "Press Shortcut keys to switch input language"
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

func (e *IBusBambooEngine) openShortcutsGUI() {
	cmd := exec.Command("/usr/lib/ibus-bamboo/keyboard-shortcut-editor", e.getShortcutString(), strconv.Itoa(e.config.DefaultInputMode))
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GTK_IM_MODULE=gtk-im-context-simple")
	out, err := cmd.Output()
	if err != nil {
		out, err = exec.Command("./keyboard-shortcut-editor", e.getShortcutString(), strconv.Itoa(e.config.DefaultInputMode)).Output()
		if err != nil {
			return
		}
	}
	if len(out) > 0 {
		e.parseShortcuts(string(out))
	} else if err != nil {
		fmt.Println("execute keyboard-shortcut-editor: ", err)
	}
}

func (e *IBusBambooEngine) parseShortcuts(s string) {
	fmt.Printf("output=(%s)\n", s)
	list := strings.Split(s, ",")
	if len(list) < len(e.config.Shortcuts) {
		return
	}
	for i := 0; i < len(e.config.Shortcuts); i++ {
		n, err := strconv.Atoi(list[i])
		if err != nil {
			fmt.Printf("ERR: failed to parse shortcut keys: %s\n", err)
		}
		e.config.Shortcuts[i] = uint32(n)
	}
}

func (e *IBusBambooEngine) getShortcutString() string {
	var s [len(e.config.Shortcuts)]string
	for i, c := range e.config.Shortcuts {
		s[i] = strconv.Itoa(int(c))
	}
	return strings.Join(s[0:], ",")
}
