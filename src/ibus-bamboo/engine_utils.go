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

func GetIBusBambooEngine() func(conn *dbus.Conn, engineName string) dbus.ObjectPath {
	objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/bamboo/%d", time.Now().UnixNano()))
	var dictionary, _ = loadDictionary(DictVietnameseCm)
	var bambooEmoji = NewBambooEmoji(DictEmojiOne)
	var mTable = NewMacroTable()
	setupConfigDir()
	go keyPressCapturing()

	return func(conn *dbus.Conn, ngName string) dbus.ObjectPath {
		var engineName = strings.ToLower(ngName)
		var config = LoadConfig(engineName)
		var inputMethod = bamboo.ParseInputMethod(config.InputMethodDefinitions, config.InputMethod)
		var preeditor = bamboo.NewEngine(inputMethod, config.Flags, dictionary)
		engine := &IBusBambooEngine{
			Engine:     ibus.BaseEngine(conn, objectPath),
			preeditor:  preeditor,
			engineName: engineName,
			config:     config,
			propList:   GetPropListByConfig(config),
			macroTable: mTable,
			dictionary: dictionary,
			emoji:      bambooEmoji,
		}
		ibus.PublishEngine(conn, objectPath, engine)

		keyPressHandler = engine.keyPressHandler

		onMouseMove = func() {
			engine.ignorePreedit = false
			x11ClipboardReset()
			engine.resetFakeBackspace()
			engine.resetBuffer()
			engine.firstTimeSendingBS = true
		}
		runtime.GC()
		debug.FreeOSMemory()
		return objectPath
	}
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
	if e.inPreeditList() || !e.inBackspaceWhiteList() {
		e.commitPreedit()
	} else {
		e.preeditor.Reset()
	}
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.preeditor.GetRawString())
}

var lookupTableConfiguration = []string{
	"Cấu hình mặc định (Pre-edit)",
	"Fix gạch chân (Surrounding Text)",
	"Fix gạch chân (Forward KeyEvent)",
	"Fix gạch chân (XTestFakeKeyEvent)",
	"Fix gạch chân (Forward as commit)",
	"Thêm vào danh sách loại trừ",
}

func (e *IBusBambooEngine) inLookupTableControlKeys(keyVal uint32) bool {
	if keyVal == IBUS_OpenLookupTable {
		return true
	}
	if idx, err := strconv.Atoi(string(keyVal)); err == nil {
		return idx < len(lookupTableConfiguration) && lookupTableConfiguration[idx] != ""
	}
	return false
}

func (e *IBusBambooEngine) openLookupTable() {
	var whiteList = [][]string{
		e.config.PreeditWhiteList,
		e.config.SurroundingTextWhiteList,
		e.config.ForwardKeyWhiteList,
		e.config.X11ClipboardWhiteList,
		e.config.DirectForwardKeyWhiteList,
		e.config.ExceptedList,
	}

	e.UpdateAuxiliaryText(ibus.NewText("Nhấn (1/2/3/4/5/6) để lưu tùy chọn của bạn"), true)

	lt := ibus.NewLookupTable()
	lt.PageSize = uint32(len(lookupTableConfiguration))
	lt.Orientation = IBUS_ORIENTATION_VERTICAL
	for i := 0; i < len(lookupTableConfiguration); i++ {
		if inStringList(whiteList[i], e.wmClasses) {
			lt.AppendLabel("*")
		} else {
			lt.AppendLabel(strconv.Itoa(i + 1))
		}
	}
	for _, ac := range lookupTableConfiguration {
		lt.AppendCandidate(ac)
	}
	e.UpdateLookupTable(lt, true)
}

func (e *IBusBambooEngine) ltProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var wmClasses = x11GetFocusWindowClass()
	e.HideLookupTable()
	fmt.Printf("keyCode 0x%04x keyval 0x%04x | %c\n", keyCode, keyVal, rune(keyVal))
	e.HideAuxiliaryText()
	if wmClasses == "" {
		return true, nil
	}
	if keyVal == IBUS_OpenLookupTable {
		return false, nil
	}
	var reset = func() {
		e.config.PreeditWhiteList = removeFromWhiteList(e.config.PreeditWhiteList, wmClasses)
		e.config.X11ClipboardWhiteList = removeFromWhiteList(e.config.X11ClipboardWhiteList, wmClasses)
		e.config.ForwardKeyWhiteList = removeFromWhiteList(e.config.ForwardKeyWhiteList, wmClasses)
		e.config.SurroundingTextWhiteList = removeFromWhiteList(e.config.SurroundingTextWhiteList, wmClasses)
		e.config.DirectForwardKeyWhiteList = removeFromWhiteList(e.config.DirectForwardKeyWhiteList, wmClasses)
		e.config.ExceptedList = removeFromWhiteList(e.config.ExceptedList, wmClasses)
	}
	switch keyVal {
	case '1':
		reset()
		e.config.PreeditWhiteList = addToWhiteList(e.config.PreeditWhiteList, wmClasses)
		break
	case '2':
		reset()
		e.config.SurroundingTextWhiteList = addToWhiteList(e.config.SurroundingTextWhiteList, wmClasses)
		break
	case '3':
		reset()
		e.config.ForwardKeyWhiteList = addToWhiteList(e.config.ForwardKeyWhiteList, wmClasses)
		break
	case '4':
		reset()
		e.config.X11ClipboardWhiteList = addToWhiteList(e.config.X11ClipboardWhiteList, wmClasses)
		break
	case '5':
		reset()
		e.config.DirectForwardKeyWhiteList = addToWhiteList(e.config.DirectForwardKeyWhiteList, wmClasses)
		break
	case '6':
		reset()
		e.config.ExceptedList = addToWhiteList(e.config.ExceptedList, wmClasses)
		break
	}

	SaveConfig(e.config, e.engineName)
	e.propList = GetPropListByConfig(e.config)
	e.RegisterProperties(e.propList)
	return true, nil
}

func (e *IBusBambooEngine) isIgnoredKey(keyVal, state uint32) bool {
	if state&IBUS_RELEASE_MASK != 0 {
		//Ignore key-up event
		return true
	}
	if e.inExceptedList() {
		if e.inLookupTableControlKeys(keyVal) {
			return false
		}
		return true
	}
	return false
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

func (e *IBusBambooEngine) canProcessKey(keyVal uint32) bool {
	if keyVal == IBUS_BackSpace {
		return true
	}
	var keyRune = rune(keyVal)
	if bamboo.IsWordBreakSymbol(keyRune) {
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
	return e.inForwardKeyList() || e.inXTestFakeKeyEventList() || e.inSurroundingTextList()
}

func (e *IBusBambooEngine) inSurroundingTextList() bool {
	return inStringList(e.config.SurroundingTextWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inDirectForwardKeyList() bool {
	return inStringList(e.config.DirectForwardKeyWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inForwardKeyList() bool {
	return e.config.IBflags&IBfakeBackspaceEnabled != 0 || inStringList(e.config.ForwardKeyWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inXTestFakeKeyEventList() bool {
	return inStringList(e.config.X11ClipboardWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inBrowserList() bool {
	return inStringList(DefaultBrowserList, e.wmClasses)
}

func (e *IBusBambooEngine) inChromeFamily() bool {
	var list = []string{
		"google-chrome:Google-chrome",
		"chromium-browser:Chromium-browser",
	}
	return inStringList(list, e.wmClasses)
}
