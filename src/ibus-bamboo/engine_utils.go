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
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var dictionary map[string]bool
var emojiTrie *TrieNode

func GetBambooEngineCreator() func(conn *dbus.Conn, engineName string) dbus.ObjectPath {
	objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/bamboo/%d", time.Now().UnixNano()))
	setupConfigDir()
	go keyPressCapturing()
	dictionary = map[string]bool{}
	emojiTrie = NewTrie()
	var engineName = strings.ToLower(EngineName)

	return func(conn *dbus.Conn, ngName string) dbus.ObjectPath {
		var engine = new(IBusBambooEngine)
		var config = LoadConfig(engineName)
		var inputMethod = bamboo.ParseInputMethod(config.InputMethodDefinitions, config.InputMethod)
		engine.Engine = ibus.BaseEngine(conn, objectPath)
		engine.engineName = engineName
		engine.preeditor = bamboo.NewEngine(inputMethod, config.Flags)
		engine.config = LoadConfig(engineName)
		engine.propList = GetPropListByConfig(config)
		ibus.PublishEngine(conn, objectPath, engine)
		go engine.init()

		return objectPath
	}
}

const KEYPRESS_DELAY_MS = 10

func (e *IBusBambooEngine) init() {
	if e.emoji == nil {
		e.emoji = NewEmojiEngine()
	}
	if e.macroTable == nil {
		e.macroTable = NewMacroTable()
		if e.config.IBflags&IBmarcoEnabled != 0 {
			e.macroTable.Enable(e.engineName)
		}
	}
	if e.config.IBflags&IBspellCheckingWithDicts != 0 {
		dictionary, _ = loadDictionary(DictVietnameseCm)
	}
	if e.config.IBflags&IBemojiDisabled == 0 {
		loadEmojiOne(DictEmojiOne)
	}
	keyPressHandler = e.keyPressHandler

	if e.config.IBflags&IBmouseCapturing != 0 {
		startMouseCapturing()
	}
	startMouseRecording()
	onMouseMove = func() {
		e.Lock()
		defer e.Unlock()
		if e.inPreeditList() || !e.inBackspaceWhiteList() {
			if e.getRawKeyLen() == 0 {
				return
			}
			e.commitPreedit(e.getPreeditString())
		}
	}
	onMouseClick = func() {
		e.Lock()
		defer e.Unlock()
		e.isFirstTimeSendingBS = true
		if e.isEmojiLTOpened {
			e.refreshEmojiCandidate()
		} else {
			e.resetFakeBackspace()
			e.resetBuffer()
			e.keyPressDelay = KEYPRESS_DELAY_MS
			if e.capabilities&IBUS_CAP_SURROUNDING_TEXT != 0 {
				//engine.ForwardKeyEvent(IBUS_Shift_R, 0, IBUS_RELEASE_MASK)
				x11SendShiftR()
				e.isSurroundingTextReady = true
				e.keyPressDelay = KEYPRESS_DELAY_MS * 10
			}
		}
	}

	runtime.GC()
	debug.FreeOSMemory()
}

var keyPressHandler = func(keyVal, keyCode, state uint32) {}
var keyPressChan = make(chan [3]uint32, 100)

func keyPressCapturing() {
	for {
		select {
		case keyEvents := <-keyPressChan:
			var keyVal, keyCode, state = keyEvents[0], keyEvents[1], keyEvents[2]
			keyPressHandler(keyVal, keyCode, state)
		}
	}
}

func (e *IBusBambooEngine) resetBuffer() {
	if e.getRawKeyLen() == 0 {
		return
	}
	if e.inPreeditList() || !e.inBackspaceWhiteList() {
		e.commitPreedit(e.getPreeditString())
	} else {
		e.preeditor.Reset()
	}
}

func (e *IBusBambooEngine) processShiftKey(keyVal, state uint32) bool {
	if keyVal == IBUS_Shift_L || keyVal == IBUS_Shift_R {
		// when press one Shift key
		if state&IBUS_SHIFT_MASK != 0 && state&IBUS_RELEASE_MASK != 0 &&
			e.config.IBflags&IBimQuickSwitchEnabled != 0 && !e.lastKeyWithShift {
			e.englishMode = !e.englishMode
			notify(e.englishMode)
			e.resetBuffer()
		}
		return true
	}
	return false
}

func (e *IBusBambooEngine) updateLastKeyWithShift(keyVal, state uint32) {
	if e.canProcessKey(keyVal, state) {
		e.lastKeyWithShift = state&IBUS_SHIFT_MASK != 0
	} else {
		e.lastKeyWithShift = false
	}
}

func (e *IBusBambooEngine) isIgnoredKey(keyVal, state uint32) bool {
	if state&IBUS_RELEASE_MASK != 0 {
		//Ignore key-up event
		return true
	}
	if keyVal == IBUS_Caps_Lock {
		return true
	}
	if e.inExceptedList() {
		if e.isInputModeLTOpened || keyVal == IBUS_OpenLookupTable {
			return false
		}
		return true
	}
	return false
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.preeditor.GetRawString())
}

func (e *IBusBambooEngine) openLookupTable() {
	var whiteList = [][]string{
		e.config.PreeditWhiteList,
		e.config.SurroundingTextWhiteList,
		e.config.ForwardKeyWhiteList,
		e.config.SLForwardKeyWhiteList,
		e.config.X11ClipboardWhiteList,
		e.config.DirectForwardKeyWhiteList,
		e.config.ExceptedList,
	}
	fmt.Println("x22 forcus", x11GetFocusWindowClass())
	var wmClasses = strings.Split(e.wmClasses, ":")
	var wmClass = e.wmClasses
	if len(wmClasses) == 2 {
		wmClass = wmClasses[1]
	}

	var lookupTableConfiguration = []string{
		"Cấu hình mặc định (Pre-edit)",
		"Sửa lỗi gạch chân (Surrounding Text)",
		"Sửa lỗi gạch chân (ForwardKeyEvent I)",
		"Sửa lỗi gạch chân (ForwardKeyEvent II)",
		"Sửa lỗi gạch chân (XTestFakeKeyEvent)",
		"Sửa lỗi gạch chân (Forward as commit)",
		"Thêm vào danh sách loại trừ (" + wmClass + ")",
	}

	e.UpdateAuxiliaryText(ibus.NewText("Nhấn (1/2/3/4/5/6/7) để lưu tùy chọn của bạn"), true)

	lt := ibus.NewLookupTable()
	lt.PageSize = uint32(len(lookupTableConfiguration))
	lt.Orientation = IBUS_ORIENTATION_VERTICAL
	var cursorPos = 0
	for i := 0; i < len(lookupTableConfiguration); i++ {
		if inStringList(whiteList[i], e.wmClasses) {
			lt.AppendLabel("*")
			cursorPos = i
		} else {
			lt.AppendLabel(strconv.Itoa(i + 1))
		}
	}
	for _, ac := range lookupTableConfiguration {
		lt.AppendCandidate(ac)
	}
	lt.SetCursorPos(uint32(cursorPos))
	e.inputModeLookupTable = lt
	e.UpdateLookupTable(lt, true)
}

func (e *IBusBambooEngine) ltProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var wmClasses = x11GetFocusWindowClass()
	//e.HideLookupTable()
	fmt.Printf("keyCode 0x%04x keyval 0x%04x | %c\n", keyCode, keyVal, rune(keyVal))
	//e.HideAuxiliaryText()
	if wmClasses == "" {
		return true, nil
	}
	if keyVal == IBUS_OpenLookupTable {
		e.closeInputModeCandidates()
		return false, nil
	}
	var keyRune = rune(keyVal)
	if keyVal == IBUS_Left || keyVal == IBUS_Up {
		e.CursorUp()
		return true, nil
	} else if keyVal == IBUS_Right || keyVal == IBUS_Down {
		e.CursorDown()
		return true, nil
	} else if keyVal == IBUS_Page_Up {
		e.PageUp()
		return true, nil
	} else if keyVal == IBUS_Page_Down {
		e.PageDown()
		return true, nil
	}
	if keyVal == IBUS_Return {
		e.commitInputModeCandidate()
		e.closeInputModeCandidates()
		return true, nil
	}
	if keyRune >= '1' && keyRune <= '9' {
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
	var wmClasses = x11GetFocusWindowClass()
	var pos = e.inputModeLookupTable.CursorPos + 1
	var reset = func() {
		e.config.PreeditWhiteList = removeFromWhiteList(e.config.PreeditWhiteList, wmClasses)
		e.config.X11ClipboardWhiteList = removeFromWhiteList(e.config.X11ClipboardWhiteList, wmClasses)
		e.config.SLForwardKeyWhiteList = removeFromWhiteList(e.config.SLForwardKeyWhiteList, wmClasses)
		e.config.ForwardKeyWhiteList = removeFromWhiteList(e.config.ForwardKeyWhiteList, wmClasses)
		e.config.SurroundingTextWhiteList = removeFromWhiteList(e.config.SurroundingTextWhiteList, wmClasses)
		e.config.DirectForwardKeyWhiteList = removeFromWhiteList(e.config.DirectForwardKeyWhiteList, wmClasses)
		e.config.ExceptedList = removeFromWhiteList(e.config.ExceptedList, wmClasses)
	}
	reset()
	switch pos {
	case 1:
		e.config.PreeditWhiteList = addToWhiteList(e.config.PreeditWhiteList, wmClasses)
	case 2:
		e.config.SurroundingTextWhiteList = addToWhiteList(e.config.SurroundingTextWhiteList, wmClasses)
	case 3:
		e.config.ForwardKeyWhiteList = addToWhiteList(e.config.ForwardKeyWhiteList, wmClasses)
	case 4:
		e.config.SLForwardKeyWhiteList = addToWhiteList(e.config.SLForwardKeyWhiteList, wmClasses)
	case 5:
		e.config.X11ClipboardWhiteList = addToWhiteList(e.config.X11ClipboardWhiteList, wmClasses)
	case 6:
		e.config.DirectForwardKeyWhiteList = addToWhiteList(e.config.DirectForwardKeyWhiteList, wmClasses)
	case 7:
		e.config.ExceptedList = addToWhiteList(e.config.ExceptedList, wmClasses)
	}

	SaveConfig(e.config, e.engineName)
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
	if state&IBUS_CONTROL_MASK != 0 ||
		state&IBUS_MOD1_MASK != 0 ||
		state&IBUS_IGNORED_MASK != 0 ||
		state&IBUS_SUPER_MASK != 0 ||
		state&IBUS_HYPER_MASK != 0 ||
		state&IBUS_META_MASK != 0 {
		return false
	}
	return true
}

func (e *IBusBambooEngine) canProcessKey(keyVal, state uint32) bool {
	if keyVal == IBUS_Space || keyVal == IBUS_BackSpace {
		return true
	}
	var keyRune = rune(keyVal)
	if bamboo.IsWordBreakSymbol(keyRune) || ('0' <= keyVal && keyVal <= '9') {
		return true
	}
	return e.preeditor.CanProcessKey(keyRune)
}

func (e *IBusBambooEngine) inExceptedList() bool {
	return inStringList(e.config.ExceptedList, e.wmClasses)
}

func (e *IBusBambooEngine) inPreeditList() bool {
	return inStringList(e.config.PreeditWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inBackspaceWhiteList() bool {
	return e.inForwardKeyList() || e.inXTestFakeKeyEventList() ||
		e.inSurroundingTextList() || e.inDirectForwardKeyList() || e.inSLForwardKeyList()
}

func (e *IBusBambooEngine) inSurroundingTextList() bool {
	return e.wmClasses != "" && inStringList(e.config.SurroundingTextWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inSLForwardKeyList() bool {
	return e.wmClasses != "" && inStringList(e.config.SLForwardKeyWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inDirectForwardKeyList() bool {
	return e.wmClasses != "" && inStringList(e.config.DirectForwardKeyWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inForwardKeyList() bool {
	return e.wmClasses != "" && inStringList(e.config.ForwardKeyWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inXTestFakeKeyEventList() bool {
	return e.wmClasses != "" && inStringList(e.config.X11ClipboardWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inBrowserList() bool {
	return inStringList(DefaultBrowserList, e.wmClasses)
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
