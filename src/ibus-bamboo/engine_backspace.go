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
	"time"
)

func (e *IBusBambooEngine) backspaceProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	if keyVal == IBUS_BackSpace {
		if e.config.IBflags&IBautoNonVnRestore == 0 {
			if e.getRawKeyLen() > 0 {
				e.preeditor.RemoveLastChar()
			}
			return false, nil
		}
		if e.nFakeBackSpace == nFakeBackspaceDefault { // just a normal backspace
			if e.getRawKeyLen() > 0 {
				oldRunes := []rune(e.getPreeditString())
				e.preeditor.RemoveLastChar()
				newRunes := []rune(e.getPreeditString())
				if len(oldRunes) == 0 {
					return false, nil
				}
				e.updatePreviousText(newRunes, oldRunes, state)
				return true, nil
			}
		} else {
			e.nFakeBackSpace--
		}
		return false, nil
	}

	if keyVal == IBUS_Return || keyVal == IBUS_KP_Enter {
		e.resetFakeBackspace()
		e.preeditor.Reset()
		return false, nil
	}

	if keyVal == IBUS_Escape {
		e.resetFakeBackspace()
		if e.getRawKeyLen() > 0 {
			e.preeditor.Reset()
			return true, nil
		}
		return false, nil
	}
	var keyRune = rune(keyVal)

	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		oldRunes := []rune(e.getPreeditString())
		e.preeditor.ProcessKey(keyRune, e.getMode())
		newRunes := []rune(e.getPreeditString())
		e.updatePreviousText(newRunes, oldRunes, state)
		return true, nil
	} else if keyVal == IBUS_space || bamboo.IsWordBreakSymbol(keyRune) {
		// macro processing
		var processedStr = e.preeditor.GetProcessedString(bamboo.VietnameseMode, false)
		if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(processedStr) {
			macText := e.macroTable.GetText(processedStr)
			macText = macText + string(keyRune)
			e.updatePreviousText([]rune(macText), []rune(processedStr), state)
			e.preeditor.Reset()
			return true, nil
		} else if e.mustFallbackToEnglish() && !e.inX11ClipboardList() {
			oldRunes := []rune(e.getPreeditString())
			newRunes := []rune(e.getComposedString())
			e.updatePreviousText(newRunes, oldRunes, state)
			e.preeditor.Reset()
		}
		e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		return false, nil
	}
	e.preeditor.Reset()
	return false, nil
}

func (e *IBusBambooEngine) updatePreviousText(newRunes, oldRunes []rune, state uint32) {
	oldLen := len(oldRunes)
	newLen := len(newRunes)
	minLen := oldLen
	if newLen < minLen {
		minLen = newLen
	}

	sameTo := -1
	for i := 0; i < minLen; i++ {
		if oldRunes[i] == newRunes[i] {
			sameTo = i
		} else {
			break
		}
	}
	diffFrom := sameTo + 1
	fmt.Println(string(oldRunes), string(newRunes), diffFrom)

	nBackSpace := 0

	if diffFrom < oldLen {
		nBackSpace += oldLen - diffFrom
	}

	e.sendBackspaceAndNewRunes(state, nBackSpace, newRunes[diffFrom:])
	mouseCaptureUnlock()
}

func (e *IBusBambooEngine) sendBackspaceAndNewRunes(state uint32, nBackSpace int, newRunes []rune) {
	if nBackSpace > 0 {
		e.SendBackSpace(state, nBackSpace)
	}
	if e.inX11ClipboardList() {
		if nBackSpace > 0 {
			e.nFakeBackSpace = nBackSpace
			x11Copy(string(newRunes))
			x11Paste(e.shortcutKeysID)
		} else {
			e.SendText(newRunes)
		}
	} else {
		e.SendText(newRunes)
	}
}

func (e *IBusBambooEngine) SendBackSpace(state uint32, n int) {
	fmt.Printf("Sendding %d backSpace.", n)

	if e.inSurroundingTextList() {
		fmt.Println("Send backspace via SurroundingText")
		e.DeleteSurroundingText(-int32(n), uint32(n))
	} else if e.inX11ClipboardList() {
		fmt.Println("Send backspace via XTestFakeKeyEvent")
		x11SendBackspace(uint32(n))
	} else if e.inForwardKeyList() {
		fmt.Println("Send backspace via IBus ForwardKeyEvent")
		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, 14, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, 14, IBUS_RELEASE_MASK)
		}
		time.Sleep(5 * time.Millisecond)
	} else {
		fmt.Println("There's something wrong with wmClasses")
	}
}

func (e *IBusBambooEngine) resetFakeBackspace() {
	e.nFakeBackSpace = nFakeBackspaceDefault
}
func (e *IBusBambooEngine) SendText(rs []rune) {
	fmt.Println("Send text", string(rs))
	e.CommitText(ibus.NewText(string(rs)))
}
