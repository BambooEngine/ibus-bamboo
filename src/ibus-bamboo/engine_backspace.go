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
	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
	"log"
)

func (e *IBusBambooEngine) backspaceProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var rawKeyLen = e.getRawKeyLen()
	if keyVal == IBUS_BackSpace {
		if rawKeyLen > 0 {
			oldRunes := []rune(e.preediter.GetProcessedString(bamboo.VietnameseMode))
			e.preediter.RemoveLastChar()
			newRunes := []rune(e.preediter.GetProcessedString(bamboo.VietnameseMode))
			e.bsUpdatePreedit(newRunes, oldRunes, state)
			return true, nil
		}

		//No thing left, just ignore
		return false, nil
	}

	if keyVal == IBUS_Return || keyVal == IBUS_KP_Enter {
		if rawKeyLen > 0 {
			e.bsCommitPreedit(keyVal)
			if e.capSurrounding() {
				return false, nil
			}
			e.ForwardKeyEvent(keyVal, keyCode, state)
			return true, nil
		} else {
			return false, nil
		}
	}

	if keyVal == IBUS_Escape {
		if rawKeyLen > 0 {
			e.bsCommitPreedit(keyVal)
			return true, nil
		}
	}

	if keyVal == IBUS_space || keyVal == IBUS_KP_Space {
		if rawKeyLen > 0 {
			e.bsCommitPreedit(0)
			if e.capSurrounding() {
				return false, nil
			}
			e.ForwardKeyEvent(keyVal, keyCode, state)
			return true, nil
		}
	}

	if (keyVal >= 'a' && keyVal <= 'z') ||
		(keyVal >= 'A' && keyVal <= 'Z') ||
		(keyVal >= '0' && keyVal <= '9') ||
		(inKeyMap(e.preediter.GetInputMethod().Keys, rune(keyVal))) {
		var keyRune = rune(keyVal)
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		if e.config.IBflags&IBautoNonVnRestore == 0 {
			oldRunes := []rune(e.preediter.GetProcessedString(bamboo.VietnameseMode))
			e.preediter.ProcessChar(keyRune, bamboo.VietnameseMode)
			newRunes := []rune(e.preediter.GetProcessedString(bamboo.VietnameseMode))
			e.bsUpdatePreedit(newRunes, oldRunes, state)
			return true, nil
		}
		oldRunes := []rune(e.getPreeditString())
		e.preediter.ProcessChar(keyRune, e.getMode())
		newRunes := []rune(e.getPreeditString())
		e.bsUpdatePreedit(newRunes, oldRunes, state)
		return true, nil
	} else {
		if rawKeyLen > 0 {
			if e.bsCommitPreedit(keyVal) {
				//lastKey already appended to commit string
				return true, nil
			} else {
				//forward lastKey
				if e.capSurrounding() {
					return false, nil
				}
				e.ForwardKeyEvent(keyVal, keyCode, state)
				return true, nil
			}
		}
		//pre-edit empty, just forward key
		return false, nil
	}
	return false, nil
}

func (e *IBusBambooEngine) capSurrounding() bool {
	return inWhiteList(e.config.SurroundingWhiteList, e.wmClasses)
}

func (e *IBusBambooEngine) bsUpdatePreedit(newRunes, oldRunes []rune, state uint32) {
	mouseCaptureUnlock()
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

	log.Println(string(oldRunes))
	log.Println(string(newRunes))
	log.Println(diffFrom)

	nBackSpace := 0
	if diffFrom < newLen && diffFrom < oldLen {
		e.SendText([]rune{0x200A}) // https://en.wikipedia.org/wiki/Whitespace_character
		nBackSpace += 1
	}

	if diffFrom < oldLen {
		nBackSpace += oldLen - diffFrom
	}

	e.SendBackSpace(state, nBackSpace)

	e.SendText(newRunes[diffFrom:])
}

func (e *IBusBambooEngine) bsCommitPreedit(lastKey uint32) bool {
	e.preediter.Reset()
	return false
}

func (e *IBusBambooEngine) SendBackSpace(state uint32, n int) {
	log.Printf("Sendding %d backSpace\n", n)
	if n == 0 {
		return
	}

	if inWhiteList(e.config.SurroundingWhiteList, e.wmClasses) {
		log.Println("Send backspace via SurroundingText")
		e.DeleteSurroundingText(-int32(n), uint32(n))
	} else if inWhiteList(e.config.X11BackspaceWhiteList, e.wmClasses) {
		log.Println("Send backspace via X11 ForwardKeyEvent")
		x11Sync(e.display)
		for i := 0; i < n; i++ {
			x11Sync(e.display)
			x11Backspace()
		}
	} else if inWhiteList(e.config.IBusBackspaceWhiteList, e.wmClasses) {
		log.Println("Send backspace via IBus ForwardKeyEvent")
		x11Flush(e.display)
		x11Sync(e.display)
		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, 14, state)
			e.ForwardKeyEvent(IBUS_BackSpace, 14, state|IBUS_RELEASE_MASK)
		}
	} else {
		log.Println("There's something wrong with wmClasses")
	}
}

func (e *IBusBambooEngine) SendText(rs []rune) {
	log.Println("Send key", string(rs))
	e.HidePreeditText()

	x11Sync(e.display)
	e.CommitText(ibus.NewText(string(rs)))
}

func (e *IBusBambooEngine) inBackspaceWhiteList(wmClasses []string) bool {
	if inWhiteList(e.config.IBusBackspaceWhiteList, wmClasses) {
		return true
	}
	if inWhiteList(e.config.X11BackspaceWhiteList, wmClasses) {
		return true
	}
	if inWhiteList(e.config.SurroundingWhiteList, wmClasses) {
		return true
	}
	return false
}
