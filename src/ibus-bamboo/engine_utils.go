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
	"time"
)

func IBusBambooEngineCreator(conn *dbus.Conn, engineName string) dbus.ObjectPath {
	objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/bamboo/%d", time.Now().UnixNano()))

	var config = LoadConfig(engineName)
	var dictionary, _ = loadDictionary(DictVietnameseCm)

	engine := &IBusBambooEngine{
		Engine:     ibus.BaseEngine(conn, objectPath),
		preediter:  bamboo.NewEngine(config.InputMethod, config.Flags, dictionary),
		engineName: engineName,
		config:     config,
		propList:   GetPropListByConfig(config),
		macroTable: NewMacroTable(),
		dictionary: dictionary,
	}
	ibus.PublishEngine(conn, objectPath, engine)

	if config.IBflags&IBmarcoEnabled != 0 {
		engine.macroTable.Enable()
	}
	go engine.startAutoCommit()

	onMouseMove = func() {
		if engine.config.IBflags&IBautoCommitWithMouseMovement == 0 {
			return
		}
		engine.ignorePreedit = false
		if engine.getRawKeyLen() == 0 {
			return
		}
		log.Println("vao day k")
		if engine.inBackspaceWhiteList(engine.wmClasses) {
			log.Println("vao nhe")
			engine.bsCommitPreedit(0)
		} else {
			engine.commitPreedit(0)
		}
	}

	return objectPath
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.getProcessedString(bamboo.EnglishMode))
}

var lookupTableControlKeys = map[uint32]string{
	'0': "Cấu hình mặc định",
	'1': "Tắt gạch chân (Surrounding Text)",
	'2': "Tắt gạch chân (IBus)",
	'3': "Tắt gạch chân (X11)",
	'4': "Thêm vào danh sách loại trừ",
}

func (e *IBusBambooEngine) inLookupTableControlKeys(keyVal uint32) bool {
	return keyVal == IBUS_OpenLookupTable || lookupTableControlKeys[keyVal] != ""
}

func (e *IBusBambooEngine) openLookupTable() {

	e.UpdateAuxiliaryText(ibus.NewText("Nhấn (1/2/3/4) để lưu tùy chọn của bạn"), true)

	lt := ibus.NewLookupTable()
	lt.AppendCandidate("Cấu hình mặc định")
	lt.AppendCandidate("Tắt gạch chân (Surrounding Text)")
	lt.AppendCandidate("Tắt gạch chân (IBus)")
	lt.AppendCandidate("Tắt gạch chân (X11)")
	lt.AppendCandidate("Thêm vào danh sách loại trừ")

	lt.AppendLabel("0:")
	lt.AppendLabel("1:")
	lt.AppendLabel("2:")
	lt.AppendLabel("3:")
	lt.AppendLabel("4:")

	e.UpdateLookupTable(lt, true)
}

func (e *IBusBambooEngine) ltProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var wmClasses = x11GetFocusWindowClass(e.display)
	e.HideLookupTable()
	fmt.Printf("keyCode 0x%04x keyval 0x%04x | %c\n", keyCode, keyVal, rune(keyVal))
	e.HideAuxiliaryText()
	var reset = func() {
		e.config.X11BackspaceWhiteList = removeWhiteList(e.config.X11BackspaceWhiteList, wmClasses)
		e.config.IBusBackspaceWhiteList = removeWhiteList(e.config.IBusBackspaceWhiteList, wmClasses)
		e.config.SurroundingWhiteList = removeWhiteList(e.config.SurroundingWhiteList, wmClasses)
		e.config.ExceptWhiteList = removeWhiteList(e.config.ExceptWhiteList, wmClasses)
	}
	switch keyVal {
	case '0':
		reset()
		break
	case '1':
		reset()
		e.config.SurroundingWhiteList = addWhiteList(e.config.SurroundingWhiteList, wmClasses)
		break
	case '2':
		reset()
		e.config.IBusBackspaceWhiteList = addWhiteList(e.config.IBusBackspaceWhiteList, wmClasses)
		break
	case '3':
		reset()
		e.config.X11BackspaceWhiteList = addWhiteList(e.config.X11BackspaceWhiteList, wmClasses)
		break
	case '4':
		e.config.ExceptWhiteList = addWhiteList(e.config.ExceptWhiteList, wmClasses)
		break
	case IBUS_OpenLookupTable:
		return false, nil
	}

	SaveConfig(e.config, e.engineName)
	e.propList = GetPropListByConfig(e.config)
	e.RegisterProperties(e.propList)
	return true, nil
}
