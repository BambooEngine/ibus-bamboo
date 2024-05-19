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
	"ibus-bamboo/config"
	"ibus-bamboo/ui"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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
		fmt.Printf("Got engine name: %s", engineName)
		var cfg = config.LoadConfig(engineName)
		var objectPath = dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/%s/%d", engineName, time.Now().UnixNano()))
		var inputMethod = bamboo.ParseInputMethod(cfg.InputMethodDefinitions, cfg.InputMethod)
		baseEngine := ibus.BaseEngine(conn, objectPath)
		var engine = NewIbusBambooEngine(engineName, config.LoadConfig(engineName), &baseEngine, bamboo.NewEngine(inputMethod, cfg.Flags))
		engine.propList = GetPropListByConfig(cfg)
		engine.shouldEnqueuKeyStrokes = true
		ibus.PublishEngine(conn, objectPath, engine)
		if *gui {
			ui.OpenGUI(engine.engineName)
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
	initConfigFiles(e.engineName)
	e.emoji = NewEmojiEngine()
	if e.macroTable == nil {
		e.macroTable = NewMacroTable(e.config.IBflags&config.IBautoCapitalizeMacro != 0)
		if e.config.IBflags&config.IBmacroEnabled != 0 {
			e.macroTable.Enable(e.engineName)
		}
	}
	keyPressHandler = e.keyPressForwardHandler

	if e.config.IBflags&config.IBmouseCapturing != 0 {
		startMouseCapturing()
		startMouseRecording()
	}
	var mouseMutex sync.Mutex
	onMouseMove = func() {
		mouseMutex.Lock()
		defer mouseMutex.Unlock()
		if e.checkInputMode(config.PreeditIM) {
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

func initConfigFiles(engineName string) {
	if sta, err := os.Stat(config.GetConfigDir(engineName)); err != nil || !sta.IsDir() {
		err = os.Mkdir(config.GetConfigDir(engineName), 0777)
		if err != nil {
			panic(err)
		}
	}
	macroPath := config.GetMacroPath(engineName)
	if _, err := os.Stat(macroPath); os.IsNotExist(err) {
		sampleFile := getEngineSubFile(sampleMactabFile)
		sample, err := ioutil.ReadFile(sampleFile)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(macroPath, sample, 0644)
		if err != nil {
			panic(err)
		}
	}
}

var keyPressHandler = func(keyVal, keyCode, state uint32) {}
var keyPressChan = make(chan [3]uint32, 100)
var lenKeyChan int32

func keyPressCapturing() {
	for keyEvents := range keyPressChan {
		atomic.StoreInt32(&lenKeyChan, int32(len(keyPressChan)))

		var keyVal, keyCode, state = keyEvents[0], keyEvents[1], keyEvents[2]
		keyPressHandler(keyVal, keyCode, state)

		atomic.AddInt32(&lenKeyChan, -1)
	}
}

var sleep = func() {
	var i = 0
	for i < 10 && atomic.LoadInt32(&lenKeyChan) > 0 {
		i++
		time.Sleep(5 * time.Millisecond)
	}
}

func (e *IBusBambooEngine) resetBuffer() {
	if e.getRawKeyLen() == 0 {
		return
	}
	if e.checkInputMode(config.PreeditIM) {
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

	if e.config.DefaultInputMode == config.UsIM {
		return true, false
	}
	if e.isShortcutKeyPressed(keyVal, state, KSRestoreKeyStrokes) {
		// fmt.Println("===== Process restoring key strokes")
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
	if e.checkInputMode(config.UsIM) {
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
		if im, ok := e.config.InputModeMapping[e.getWmClass()]; ok && config.ImLookupTable[im] != "" {
			return im
		}
	}
	if _, ok := config.ImLookupTable[e.config.DefaultInputMode]; ok {
		return e.config.DefaultInputMode
	}
	return config.PreeditIM
}

func (e *IBusBambooEngine) openLookupTable() {
	var wmClasses = strings.Split(e.getWmClass(), ":")
	var wmClass = e.getWmClass()
	if len(wmClasses) == 2 {
		wmClass = wmClasses[1]
	}

	e.UpdateAuxiliaryText(ibus.NewText("Nhấn (1/2/3/4/5/6/7) để lưu tùy chọn của bạn"), true)

	lt := ibus.NewLookupTable()
	lt.PageSize = uint32(len(config.ImLookupTable))
	lt.Orientation = IBusOrientationVertical
	for im := 1; im <= len(config.ImLookupTable); im++ {
		if e.getInputMode() == im {
			lt.AppendLabel("*")
			lt.SetCursorPos(uint32(im - 1))
		} else {
			lt.AppendLabel(strconv.Itoa(im))
		}
		if im == config.UsIM {
			lt.AppendCandidate(config.ImLookupTable[im] + " (" + wmClass + ")")
		} else {
			lt.AppendCandidate(config.ImLookupTable[im])
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

	config.SaveConfig(e.config, e.engineName)
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

func isValidState(state uint32) bool {
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

func (e *IBusBambooEngine) isPrintableKey(state, keyVal uint32) bool {
	return isValidState(state) && e.isValidKeyVal(keyVal)
}

func (e *IBusBambooEngine) getCommitText(keyVal, keyCode, state uint32) (newText string, IsWordBreakSymbol bool) {
	var keyRune = rune(keyVal)
	isPrintableKey := e.isPrintableKey(state, keyVal)
	oldText := e.getPreeditString()
	// restore key strokes by pressing Shift + Space
	if e.shouldRestoreKeyStrokes {
		e.shouldRestoreKeyStrokes = false
		e.preeditor.RestoreLastWord(!bamboo.HasAnyVietnameseRune(oldText))
		return e.getPreeditString(), false
	}
	var keyS string
	if isPrintableKey {
		keyS = string(keyRune)
	}
	if isPrintableKey && e.preeditor.CanProcessKey(keyRune) {
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
				// TODO: THIS IS A HACK
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
		} else if e.config.IBflags&config.IBmacroEnabled != 0 {
			return e.getProcessedString(bamboo.PunctuationMode), false
		} else {
			return e.getPreeditString(), false
		}
	} else if e.config.IBflags&config.IBmacroEnabled != 0 {
		// macro processing
		if isPrintableKey && e.macroTable.HasPrefix(oldText+keyS) {
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return oldText + keyS, false
		}
		if e.macroTable.HasKey(oldText) {
			if isPrintableKey {
				return e.expandMacro(oldText) + keyS, true
			}
			return e.expandMacro(oldText), true
		}
	}
	return e.handleNonVnWord(keyVal, keyCode, state), true
}

func (e *IBusBambooEngine) handleNonVnWord(keyVal, keyCode, state uint32) string {
	var (
		keyS           string
		keyRune        = rune(keyVal)
		isPrintableKey = e.isPrintableKey(state, keyVal)
		oldText        = e.getPreeditString()
	)
	if isPrintableKey {
		keyS = string(keyRune)
	}
	if bamboo.HasAnyVietnameseRune(oldText) && e.mustFallbackToEnglish() {
		e.preeditor.RestoreLastWord(false)
		newText := e.preeditor.GetProcessedString(bamboo.PunctuationMode|bamboo.EnglishMode) + keyS
		if isPrintableKey {
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		}
		return newText
	}
	if isPrintableKey {
		e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		return oldText + keyS
	}
	// Ctrl + A is treasted as a WBS
	return oldText + keyS
}

func (e *IBusBambooEngine) getMacroText() (bool, string) {
	if e.config.IBflags&config.IBmacroEnabled == 0 {
		return false, ""
	}
	var text = e.preeditor.GetProcessedString(bamboo.PunctuationMode)
	if e.macroTable.HasKey(text) {
		return true, e.expandMacro(text)
	}
	return false, ""
}

func (e *IBusBambooEngine) getFakeBackspace() int32 {
	return atomic.LoadInt32(&e.nFakeBackSpace)
}

func (e *IBusBambooEngine) setFakeBackspace(n int32) {
	atomic.StoreInt32(&e.nFakeBackSpace, n)
}

func (e *IBusBambooEngine) addFakeBackspace(n int32) {
	atomic.AddInt32(&e.nFakeBackSpace, n)
}

func (e *IBusBambooEngine) isValidKeyVal(keyVal uint32) bool {
	var keyRune = rune(keyVal)
	if keyVal == IBusBackSpace || bamboo.IsWordBreakSymbol(keyRune) {
		return true
	}
	if ok, _ := e.getMacroText(); ok && keyVal == IBusTab {
		return true
	}
	return e.preeditor.CanProcessKey(keyRune)
}

func (e *IBusBambooEngine) inBackspaceWhiteList() bool {
	var inputMode = e.getInputMode()
	for _, im := range config.ImBackspaceList {
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
