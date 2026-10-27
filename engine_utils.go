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
	"ibus-bamboo/emoji"
	"ibus-bamboo/mactab"
	"ibus-bamboo/ui"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/BambooEngine/bamboo-core"
	ibus "github.com/BambooEngine/goibus"
	"github.com/godbus/dbus/v5"
)

var dictionary = map[string]bool{}

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
		engineRef = engine
		engine.PropList = GetPropListByConfig(cfg)
		engine.ShouldEnqueuKeyStrokes = true
		ibus.PublishEngine(conn, objectPath, engine)
		if *gui {
			ui.OpenGUI(engine.EngineName)
			os.Exit(0)
		}
		go engine.init()

		return objectPath
	}
}

const KeypressDelayMs = 10

func (e *IBusBambooEngine) isShortcutKeyEnable(ski uint) bool {
	if int(ski+2) > len(e.Config.Shortcuts) {
		return false
	}
	l := e.Config.Shortcuts[ski : ski+2]
	return l[1] > 0
}

func (e *IBusBambooEngine) init() {
	initConfigFiles(e.EngineName)
	e.Emoji = emoji.NewEmojiEngine()
	if e.MacroTable == nil {
		e.MacroTable = mactab.NewMacroTable(e.Config.IBflags&config.IBautoCapitalizeMacro != 0)
		if e.Config.IBflags&config.IBmacroEnabled != 0 {
			mactab.Enable(e.MacroTable, e.EngineName)
		}
	}
	keyPressHandler = e.keyPressForwardHandler
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
		e.Preeditor.Reset()
	}
}

func (e *IBusBambooEngine) checkWmClass(newId string) {
	if e.WmClasses != newId {
		e.WmClasses = newId
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
	shortcuts := e.Config.Shortcuts[shortcut : shortcut+2]
	ret := shortcuts[0] == realState && shortcuts[1] == lowerKey
	// fmt.Println("...isShortcutKeyPressed=", ret, ret && !e.LastKeyWithShift, shortcuts)
	if realState == 1 && shortcut == KSViEnSwitch {
		return ret && !e.LastKeyWithShift
	}
	return ret
}

func (e *IBusBambooEngine) processShortcutKey(keyVal, keyCode, state uint32) (bool, bool) {
	if keyVal == IBusCapsLock {
		return true, false
	}
	// fmt.Println("===Process shortcut for emoji selector")
	if e.isShortcutKeyPressed(keyVal, state, KSEmojiDialog) &&
		!e.IsEmojiLTOpened {
		e.resetBuffer()
		e.IsEmojiLTOpened = true
		e.LastKeyWithShift = true
		e.openEmojiList()
		return true, true
	}
	if e.IsEmojiLTOpened {
		return true, e.emojiProcessKeyEvent(keyVal, keyCode, state)
	}
	// fmt.Println("====== Process hexadecimal key pressed")
	if e.isShortcutKeyPressed(keyVal, state, KSHexadecimal) {
		e.resetBuffer()
		e.IsInHexadecimal = true
		e.setupHexadecimalProcessKeyEvent()
		return true, true
	}
	if e.IsInHexadecimal {
		if e.isShortcutKeyPressed(keyVal, state, KSHexadecimal) {
			e.closeHexadecimalInput()
			e.updateLastKeyWithShift(keyVal, state)
			return true, false
		}
		return true, e.hexadecimalProcessKeyEvent(keyVal, keyCode, state)
	}

	if e.Config.DefaultInputMode == config.UsIM {
		return true, false
	}
	if e.isShortcutKeyPressed(keyVal, state, KSRestoreKeyStrokes) {
		// fmt.Println("===== Process restoring key strokes")
		e.ShouldRestoreKeyStrokes = true
		return false, false
	}
	// fmt.Println("===Process shortcut for input method switcher")
	if e.isShortcutKeyPressed(keyVal, state, KSViEnSwitch) {
		e.EnglishMode = !e.EnglishMode
		notify(e.EnglishMode)
		e.resetBuffer()
		return true, true
	}
	// fmt.Println("====== Process shortcut for input mode switch")
	if e.IsInputModeLTOpened {
		return e.ltProcessKeyEvent(keyVal, keyCode, state)
	} else if e.isShortcutKeyPressed(keyVal, state, KSInputModeSwitch) &&
		e.getWmClass() != "" {
		e.resetBuffer()
		e.IsInputModeLTOpened = true
		e.LastKeyWithShift = true
		e.openLookupTable()
		return true, true
	}

	if keyVal == IBusShiftL || keyVal == IBusShiftR {
		return true, false
	}
	if e.checkInputMode(config.UsIM) {
		if e.IsInputModeLTOpened && e.isShortcutKeyPressed(keyVal, state, KSInputModeSwitch) {
			return false, false
		}
		return true, false
	}
	if e.EnglishMode {
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

	if upperSpecialKey, found := keyMapping[keyRune]; found && inKeyList(e.Preeditor.GetInputMethod().AppendingKeys, keyRune) {
		keyRune = upperSpecialKey
	}
	return keyRune
}

func (e *IBusBambooEngine) updateLastKeyWithShift(keyVal, state uint32) {
	if e.Preeditor.CanProcessKey(rune(keyVal)) {
		e.LastKeyWithShift = state&IBusShiftMask != 0
	} else {
		e.LastKeyWithShift = false
	}
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.Preeditor.GetProcessedString(bamboo.EnglishMode | bamboo.FullText))
}

func (e *IBusBambooEngine) runeCount() int {
	return utf8.RuneCountInString(e.getPreeditString())
}

func (e *IBusBambooEngine) getInputMode() int {
	if e.getWmClass() != "" {
		if im, ok := e.Config.InputModeMapping[e.getWmClass()]; ok && config.ImLookupTable[im] != "" {
			return im
		}
	}
	if _, ok := config.ImLookupTable[e.Config.DefaultInputMode]; ok {
		return e.Config.DefaultInputMode
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
	e.InputModeLookupTable = lt
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
			if e.InputModeLookupTable.SetCursorPos(uint32(pos - 1)) {
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
	var im = e.InputModeLookupTable.CursorPos + 1
	e.Config.InputModeMapping[e.getWmClass()] = int(im)

	config.SaveConfig(e.Config, e.EngineName)
	e.PropList = GetPropListByConfig(e.Config)
	e.RegisterProperties(e.PropList)
}

func (e *IBusBambooEngine) closeInputModeCandidates() {
	e.InputModeLookupTable = nil
	e.UpdateLookupTable(ibus.NewLookupTable(), true) // workaround for issue #18
	e.HidePreeditText()
	e.HideLookupTable()
	e.HideAuxiliaryText()
	e.IsInputModeLTOpened = false
}

func (e *IBusBambooEngine) updateInputModeLT() {
	var visible = len(e.InputModeLookupTable.Candidates) > 0
	e.UpdateLookupTable(e.InputModeLookupTable, visible)
}

func isValidState(state uint32) bool {
	if state&IBusControlMask != 0 ||
		state&IBusMod1Mask != 0 ||
		state&IBusMod4Mask != 0 ||
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
	if e.ShouldRestoreKeyStrokes {
		e.ShouldRestoreKeyStrokes = false
		e.Preeditor.RestoreLastWord(!bamboo.HasAnyVietnameseRune(oldText))
		return e.getPreeditString(), false
	}
	var keyS string
	if isPrintableKey {
		keyS = string(keyRune)
	}
	if isPrintableKey && e.Preeditor.CanProcessKey(keyRune) {
		if state&IBusLockMask != 0 {
			keyRune = e.toUpper(keyRune)
		}
		e.Preeditor.ProcessKey(keyRune, e.getBambooInputMode())
		if inKeyList(e.Preeditor.GetInputMethod().AppendingKeys, keyRune) {
			var newText string
			if e.shouldFallbackToEnglish(true) {
				newText = e.getProcessedString(bamboo.EnglishMode)
			} else {
				newText = e.getProcessedString(bamboo.VietnameseMode)
			}
			if fullSeq := e.Preeditor.GetProcessedString(bamboo.VietnameseMode); len(fullSeq) > 0 && rune(fullSeq[len(fullSeq)-1]) == keyRune {
				// [[ => [
				var ret = e.getPreeditString()
				var lastRune = rune(ret[len(ret)-1])
				var isWordBreakRune = bamboo.IsWordBreakSymbol(lastRune)
				// TODO: THIS IS A HACK
				if isWordBreakRune {
					e.Preeditor.RemoveLastChar(false)
					e.Preeditor.ProcessKey(' ', bamboo.EnglishMode)
				}
				return ret, isWordBreakRune
			} else if l := []rune(newText); len(l) > 0 && keyRune == l[len(l)-1] {
				// f] => f]
				var isWordBreakRune = bamboo.IsWordBreakSymbol(keyRune)
				if isWordBreakRune {
					e.Preeditor.RemoveLastChar(false)
					e.Preeditor.ProcessKey(' ', bamboo.EnglishMode)
				}
				return oldText + string(keyRune), isWordBreakRune
			} else {
				// ] => o?
				return e.getPreeditString(), false
			}
		} else if e.Config.IBflags&config.IBmacroEnabled != 0 {
			return e.getProcessedString(bamboo.PunctuationMode), false
		} else {
			return e.getPreeditString(), false
		}
	} else if e.Config.IBflags&config.IBmacroEnabled != 0 {
		// macro processing
		if isPrintableKey && mactab.HasPrefix(e.MacroTable, oldText+keyS) {
			e.Preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return oldText + keyS, false
		}
		if mactab.HasKey(e.MacroTable, oldText) {
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
		e.Preeditor.RestoreLastWord(false)
		newText := e.Preeditor.GetProcessedString(bamboo.PunctuationMode|bamboo.EnglishMode) + keyS
		if isPrintableKey {
			e.Preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		}
		return newText
	}
	if isPrintableKey {
		e.Preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		return oldText + keyS
	}
	// Ctrl + A is treasted as a WBS
	return oldText + keyS
}

func (e *IBusBambooEngine) getMacroText() (bool, string) {
	if e.Config.IBflags&config.IBmacroEnabled == 0 {
		return false, ""
	}
	var text = e.Preeditor.GetProcessedString(bamboo.PunctuationMode)
	if mactab.HasKey(e.MacroTable, text) {
		return true, e.expandMacro(text)
	}
	return false, ""
}

func (e *IBusBambooEngine) getFakeBackspace() int32 {
	return atomic.LoadInt32(&e.NFakeBackSpace)
}

func (e *IBusBambooEngine) setFakeBackspace(n int32) {
	atomic.StoreInt32(&e.NFakeBackSpace, n)
}

func (e *IBusBambooEngine) addFakeBackspace(n int32) {
	atomic.AddInt32(&e.NFakeBackSpace, n)
}

func (e *IBusBambooEngine) isValidKeyVal(keyVal uint32) bool {
	var keyRune = rune(keyVal)
	if keyVal == IBusBackSpace || bamboo.IsWordBreakSymbol(keyRune) {
		return true
	}
	if ok, _ := e.getMacroText(); ok && keyVal == IBusTab {
		return true
	}
	return e.Preeditor.CanProcessKey(keyRune)
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
	return e.WmClasses
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
