/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Nguyen Cong Hoang <hoangnc.jp@gmail.com>
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
	"time"
)

const (
	DiffNumpadKeypad = IBUS_KP_0 - IBUS_0
)

var (
	printableKeyCode = map[uint32]bool{
		0x0039: true,
		0x0002: true,
		0x0003: true,
		0x0004: true,
		0x0005: true,
		0x0006: true,
		0x0007: true,
		0x0008: true,
		0x0009: true,
		0x000a: true,
		0x000b: true,
		0x000c: true,
		0x000d: true,
		0x007c: true,
		0x001a: true,
		0x001b: true,
		0x0027: true,
		0x0028: true,
		0x002b: true,
		0x0033: true,
		0x0034: true,
		0x0035: true,
		0x0059: true,
	}
)

func IBusBambooEngineCreator(conn *dbus.Conn, engineName string) dbus.ObjectPath {

	objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/bamboo/%d", time.Now().UnixNano()))

	var config = LoadConfig(engineName)

	engine := &IBusBambooEngine{
		Engine:     ibus.BaseEngine(conn, objectPath),
		preediter:  bamboo.NewEngine(config.InputMethod, config.Flags),
		engineName: engineName,
		config:     config,
		propList:   GetPropListByConfig(config),
	}
	ibus.PublishEngine(conn, objectPath, engine)

	go engine.startAutoCommit()

	onMouseClick = func() {
		engine.Lock()
		defer engine.Unlock()
		var rawKeyLen = engine.getRawKeyLen()
		if rawKeyLen > 0 {
			engine.HidePreeditText()
			engine.preediter.Reset()
		}
	}
	onMouseMove = func() {
		engine.Lock()
		defer engine.Unlock()
		var rawKeyLen = engine.getRawKeyLen()
		if rawKeyLen > 0 {
			engine.commitPreedit(0)
		}
	}

	return objectPath
}

var keyPressChan = make(chan uint32)

func (e *IBusBambooEngine) startAutoCommit() {
	for {
		select {
		case <-keyPressChan:
			break
		case <-time.After(3 * time.Second):
			var rawKeyLen = e.getRawKeyLen()
			if rawKeyLen > 0 {
				e.commitPreedit(0)
			}
		}
	}
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.preediter.GetProcessedString(bamboo.EnglishMode))
}

func (e *IBusBambooEngine) updatePreedit() {
	var processedStr = e.getPreeditString()
	var preeditLen = uint32(len([]rune(processedStr)))
	if preeditLen > 0 {
		var ibusText = ibus.NewText(processedStr)
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)

		e.UpdatePreeditTextWithMode(ibusText, preeditLen, true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	} else {
		e.HidePreeditText()
		e.preediter.Reset()
	}
}

func (e *IBusBambooEngine) shouldFallbackToEnglish() bool {
	if e.config.Flags&bamboo.EspellCheckEnabled == 0 {
		return false
	}
	if e.preediter.IsSpellingCorrect(bamboo.NoTone) {
		return false
	}
	if e.preediter.IsSpellingSensible(bamboo.NoTone) {
		return false
	}
	return true
}

func (e *IBusBambooEngine) getCommitString() string {
	var processedStr string
	if e.config.Flags&bamboo.EspellCheckEnabled != 0 && !e.preediter.IsSpellingCorrect(bamboo.NoTone) {
		processedStr = e.preediter.GetProcessedString(bamboo.EnglishMode)
		return processedStr
	}
	processedStr = e.preediter.GetProcessedString(bamboo.VietnameseMode)
	processedStr = bamboo.Encode(e.config.Charset, processedStr)
	return processedStr
}

func (e *IBusBambooEngine) getPreeditString() string {
	var processedStr string
	if e.shouldFallbackToEnglish() {
		processedStr = e.preediter.GetProcessedString(bamboo.EnglishMode)
		return processedStr
	}
	processedStr = e.preediter.GetProcessedString(bamboo.VietnameseMode)
	return processedStr
}

func (e *IBusBambooEngine) commitPreedit(lastKey uint32) bool {
	var keyAppended = false
	var commitStr string
	commitStr += e.getCommitString()
	e.preediter.Reset()

	//Convert num-pad key to normal number
	if (lastKey >= IBUS_KP_0 && lastKey <= IBUS_KP_9) ||
		(lastKey >= IBUS_KP_Multiply && lastKey <= IBUS_KP_Divide) {
		lastKey = lastKey - DiffNumpadKeypad
	}

	if lastKey >= 0x20 && lastKey <= 0xFF {
		//append printable keys
		commitStr += string(lastKey)
		keyAppended = true
	}

	e.HidePreeditText()
	e.CommitText(ibus.NewText(commitStr))

	return keyAppended
}
