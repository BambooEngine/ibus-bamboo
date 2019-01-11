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

func (e *IBusBambooEngine) backspaceProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var rawKeyLen = e.getRawKeyLen(false)
	if keyVal == IBUS_BackSpace {
		log.Println("Number of fake backspaces::", e.nFakeBackSpace)
		if e.nFakeBackSpace == nFakeBackspaceDefault { // just a normal backspace
			if rawKeyLen > 0 {
				oldRunes := []rune(e.getPreeditString())
				e.preediter.RemoveLastChar()
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
		e.preediter.Reset()
		return false, nil
	}

	if keyVal == IBUS_Escape {
		e.resetFakeBackspace()
		if e.getRawKeyLen(false) > 0 {
			e.preediter.Reset()
			return true, nil
		}
		return false, nil
	}
	var keyRune = rune(keyVal)

	if (keyVal >= 'a' && keyVal <= 'z') ||
		(keyVal >= 'A' && keyVal <= 'Z') ||
		(inKeyMap(e.preediter.GetInputMethod().Keys, rune(keyVal))) {
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		if e.config.IBflags&IBautoNonVnRestore == 0 {
			oldRunes := []rune(e.preediter.GetProcessedString(bamboo.VietnameseMode, true))
			e.preediter.ProcessChar(keyRune, bamboo.VietnameseMode)
			newRunes := []rune(e.preediter.GetProcessedString(bamboo.VietnameseMode, true))
			e.updatePreviousText(newRunes, oldRunes, state)
			return true, nil
		}
		oldRunes := []rune(e.getPreeditString())
		e.preediter.ProcessChar(keyRune, e.getMode())
		newRunes := []rune(e.getPreeditString())
		e.updatePreviousText(newRunes, oldRunes, state)
		return true, nil
	} else if keyVal == IBUS_space || isWordBreakSymbol(keyRune) {
		// macro processing
		var processedStr = e.preediter.GetProcessedString(bamboo.VietnameseMode, true)
		if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(processedStr) {
			macText := e.macroTable.GetText(processedStr)
			macText = macText + string(keyRune)
			e.updatePreviousText([]rune(macText), []rune(processedStr), state)
			e.preediter.Reset()
			return true, nil
		}
		e.preediter.ProcessChar(keyRune, bamboo.EnglishMode)
		return false, nil
	}
	e.preediter.Reset()
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
	log.Println(string(oldRunes), string(newRunes), diffFrom)

	nBackSpace := 0
	if e.inSurroundingList() && diffFrom < newLen && diffFrom < oldLen {
		e.SendText([]rune{0x200A}) // https://en.wikipedia.org/wiki/Whitespace_character
		nBackSpace += 1
	}

	if diffFrom < oldLen {
		nBackSpace += oldLen - diffFrom
	}

	if nBackSpace > 0 {
		e.SendBackSpace(state, nBackSpace)
	}
	if e.inX11BackspaceList() {
		if nBackSpace > 0 {
			e.nFakeBackSpace = nBackSpace
			x11Copy(string(newRunes[diffFrom:]))
			x11Paste(e.shortcutKeysID)
		} else {
			e.SendText(newRunes[diffFrom:])
		}
	} else {
		e.SendText(newRunes[diffFrom:])
	}
	mouseCaptureUnlock()
}

func (e *IBusBambooEngine) SendBackSpace(state uint32, n int) {
	log.Printf("Sendding %d backSpace.", n)
	if n == 0 {
		return
	}

	if inWhiteList(e.config.SurroundingWhiteList, e.wmClasses) {
		fmt.Println("Send backspace via SurroundingText")
		e.DeleteSurroundingText(-int32(n), uint32(n))
	} else if inWhiteList(e.config.X11BackspaceWhiteList, e.wmClasses) {
		fmt.Println("Send backspace via X11 KeyEvent")
		x11SendBackspace(uint32(n))
	} else if inWhiteList(e.config.IBusBackspaceWhiteList, e.wmClasses) {
		fmt.Println("Send backspace via IBus ForwardKeyEvent")
		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, 14, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, 14, IBUS_RELEASE_MASK)
			time.Sleep(5 * time.Millisecond)
		}
	} else {
		fmt.Println("There's something wrong with wmClasses")
	}
}

func (e *IBusBambooEngine) resetFakeBackspace() {
	e.nFakeBackSpace = nFakeBackspaceDefault
}
func (e *IBusBambooEngine) SendText(rs []rune) {
	log.Println("Send text", string(rs))
	e.CommitText(ibus.NewText(string(rs)))
}

func (e *IBusBambooEngine) inBackspaceWhiteList(wmClasses []string) bool {
	return e.inIBusForwardList() || e.inX11BackspaceList() || e.inSurroundingList()
}

func (e *IBusBambooEngine) inSurroundingList() bool {
	return inWhiteList(e.config.SurroundingWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inIBusForwardList() bool {
	return inWhiteList(e.config.IBusBackspaceWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) inX11BackspaceList() bool {
	return inWhiteList(e.config.X11BackspaceWhiteList, e.wmClasses)
}
