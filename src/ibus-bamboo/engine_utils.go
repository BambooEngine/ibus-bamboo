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
	"strconv"
	"time"
)

func IBusBambooEngineCreator(conn *dbus.Conn, engineName string) dbus.ObjectPath {
	objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/bamboo/%d", time.Now().UnixNano()))

	var config = LoadConfig(engineName)
	var dictionary, _ = loadDictionary(DictVietnameseCm)

	engine := &IBusBambooEngine{
		Engine:         ibus.BaseEngine(conn, objectPath),
		preeditor:      bamboo.NewEngine(config.InputMethod, config.Flags, dictionary),
		engineName:     engineName,
		config:         config,
		propList:       GetPropListByConfig(config),
		macroTable:     NewMacroTable(),
		dictionary:     dictionary,
		nFakeBackSpace: nFakeBackspaceDefault,
		emoji:          NewBambooEmoji(DictEmojiOne),
	}
	ibus.PublishEngine(conn, objectPath, engine)

	if config.IBflags&IBmarcoEnabled != 0 {
		engine.macroTable.Enable()
	}
	go engine.startAutoCommit()

	onMouseMove = func() {
		engine.ignorePreedit = false
		x11ClipboardReset()
		engine.resetFakeBackspace()
		if engine.inBackspaceWhiteList() {
			engine.preeditor.Reset()
		} else if engine.isEmojiTableOpened {
			if cps := engine.emoji.Query(); len(cps) > 0 {
				engine.CommitText(ibus.NewText(cps[0]))
			}
			engine.emoji.Reset()
			engine.HidePreeditText()
			engine.HideLookupTable()
			engine.HideAuxiliaryText()
		} else {
			engine.commitPreedit(0)
		}
	}

	return objectPath
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.preeditor.GetRawString())
}

var lookupTableConfiguration = []string{
	"Cấu hình mặc định (Pre-edit)",
	"Tắt gạch chân (Surrounding Text)",
	"Tắt gạch chân (Forward key event)",
	"Tắt gạch chân (X11 Clipboard)",
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
		e.config.ExceptedWhiteList,
	}

	e.UpdateAuxiliaryText(ibus.NewText("Nhấn (0/1/2/3/4) để lưu tùy chọn của bạn"), true)

	lt := ibus.NewLookupTable()
	lt.Orientation = IBUS_ORIENTATION_VERTICAL
	for _, ac := range lookupTableConfiguration {
		lt.AppendCandidate(ac)
	}
	for lb, _ := range lookupTableConfiguration {
		if inStringList(whiteList[lb], e.wmClasses) {
			lt.AppendLabel("*")
		} else {
			lt.AppendLabel(strconv.Itoa(lb))
		}
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
		e.config.ExceptedWhiteList = removeFromWhiteList(e.config.ExceptedWhiteList, wmClasses)
	}
	switch keyVal {
	case '0':
		reset()
		e.config.PreeditWhiteList = addToWhiteList(e.config.PreeditWhiteList, wmClasses)
		break
	case '1':
		reset()
		e.config.SurroundingTextWhiteList = addToWhiteList(e.config.SurroundingTextWhiteList, wmClasses)
		break
	case '2':
		reset()
		e.config.ForwardKeyWhiteList = addToWhiteList(e.config.ForwardKeyWhiteList, wmClasses)
		break
	case '3':
		reset()
		e.config.X11ClipboardWhiteList = addToWhiteList(e.config.X11ClipboardWhiteList, wmClasses)
		break
	case '4':
		reset()
		e.config.ExceptedWhiteList = addToWhiteList(e.config.ExceptedWhiteList, wmClasses)
		break
	}

	SaveConfig(e.config, e.engineName)
	e.propList = GetPropListByConfig(e.config)
	e.RegisterProperties(e.propList)
	return true, nil
}

func (e *IBusBambooEngine) isIgnoredKey(keyVal, state uint32) bool {
	if e.inX11ClipboardList() {
		if state&IBUS_SHIFT_MASK != 0 && keyVal == IBUS_KEY_Insert {
			return true
		}
	}
	if state&IBUS_RELEASE_MASK != 0 {
		//Ignore key-up event
		return true
	}
	if keyVal == IBUS_Shift_L {
		if state&IBUS_SHIFT_MASK == 0 {
			e.shortcutKeysID = 1
		}
		return true
	} else if keyVal == IBUS_Shift_R {
		if state&IBUS_SHIFT_MASK == 0 {
			e.shortcutKeysID = 0
		}
		return true
	}
	if e.inExceptedList() {
		if e.inLookupTableControlKeys(keyVal) {
			return false
		}
		return true
	}
	return e.zeroLocation
}

func (e *IBusBambooEngine) reset() {
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) inExceptedList() bool {
	return inStringList(e.config.ExceptedWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inPreeditList() bool {
	return inStringList(e.config.PreeditWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inBackspaceWhiteList() bool {
	return e.inForwardKeyList() || e.inX11ClipboardList() || e.inSurroundingTextList()
}

func (e *IBusBambooEngine) inSurroundingTextList() bool {
	return inStringList(e.config.SurroundingTextWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inForwardKeyList() bool {
	return inStringList(e.config.ForwardKeyWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inX11ClipboardList() bool {
	return inStringList(e.config.X11ClipboardWhiteList, e.wmClasses)
}
