/*
 * Bamboo - A Vietnamese Input method editor - Hexadecimal emulator
 * for Bamboo
 * Copyright (C) 2021 Tran Duc Binh <binhtran432k@gmail.com>
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
	"log"
	"strconv"

	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

func (e *IBusBambooEngine) hexadecimalProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var rawKeyLen = e.getRawKeyLen()
	if (keyVal >= 0xffb0 && keyVal <= 0xffb9) {
		keyVal = keyVal - 0xffb0 + 0x0030
	}
	var keyRune = rune(keyVal)
	var mode = bamboo.EnglishMode | bamboo.FullText
	var oldText = e.getProcessedString(mode)
	defer e.updateLastKeyWithShift(keyVal, state)


	if rawKeyLen == 0 || oldText[0] != 'u' {
		e.closeHexadecimalInput()
		return false, nil
	}
	if keyVal == IBusEscape {
		e.closeHexadecimalInput()
		return true, nil
	}
	if keyVal == IBusBackSpace {
		if rawKeyLen > 2 {
			e.preeditor.RemoveLastChar(true)
			e.updateHexadecimal(e.getProcessedString(mode))
		} else {
			e.closeHexadecimalInput()
		}
		return true, nil
	}
	if keyVal == IBusSpace || keyVal == IBusReturn || keyVal == 0xff8d {
		if rawKeyLen > 1 {
			value, err := strconv.ParseInt(oldText[1:], 16, 64)
			if err != nil || value > 0x10ffff {
				log.Println("Input is out of range")
			} else {
				log.Printf("Commit Text [%s]\n", string(value))
				e.CommitText(ibus.NewText(string(value)))
			}
		}
		e.closeHexadecimalInput()
		return true, nil
	}

	if (keyRune >= '0' && keyRune <= '9') || (keyRune >= 'A' && keyRune <= 'F') || (keyRune >= 'a' && keyRune <= 'f') {
		if !e.isValidState(state) || !e.canProcessKey(keyVal) {
			return true, nil
		}
		if e.canProcessKey(keyVal) {
			e.preeditor.ProcessKey(keyRune, mode)
		}
		e.updateHexadecimal(e.getProcessedString(mode))
	}
	return true, nil
}

func (e *IBusBambooEngine) setupHexadecimalProcessKeyEvent() (bool, *dbus.Error) {
	var keyVal = uint32(117)
	var state = uint32(0)
	var keyRune = rune(keyVal)
	var mode = bamboo.EnglishMode | bamboo.FullText
	defer e.updateLastKeyWithShift(keyVal, state)

	if e.canProcessKey(keyVal) {
		e.preeditor.ProcessKey(keyRune, mode)
	}
	e.updateHexadecimal(e.getProcessedString(mode))

	return true, nil
}

func (e *IBusBambooEngine) closeHexadecimalInput() {
	e.HidePreeditText()
	e.preeditor.Reset()
	e.isInHexadecimal = false
}

func (e *IBusBambooEngine) updateHexadecimal(processedStr string) {
	var encodedStr = e.encodeText(processedStr)
	var preeditLen = uint32(len([]rune(encodedStr)))
	if preeditLen == 0 {
		e.HidePreeditText()
		e.CommitText(ibus.NewText(""))
		return
	}
	var ibusText = ibus.NewText(encodedStr)

	ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	e.UpdatePreeditTextWithMode(ibusText, preeditLen, true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
}
