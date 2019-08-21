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
	"github.com/godbus/dbus"
	"log"
	"time"
)

func (e *IBusBambooEngine) bsProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	if e.getRawKeyLen() == 0 {
		e.RequireSurroundingText()
	}
	if e.inXTestFakeKeyEventList() || e.inX11ShiftLeftList() || e.inSurroundingTextList() {
		// we don't want to use ForwardKeyEvent api in X11 XTestFakeKeyEvent and Surrounding Text mode
		var sleep = func() {
			for len(keyPressChan) > 0 {
				time.Sleep(5 * time.Millisecond)
			}
		}
		if keyVal == IBUS_Left && state&IBUS_SHIFT_MASK != 0 {
			if e.nFakeShiftLeft > 0 {
				e.nFakeShiftLeft--
			}
			return false, nil
		}
		if !e.isValidState(state) || !e.canProcessKey(keyVal, state) {
			e.preeditor.Reset()
			e.resetFakeBackspace()
			e.firstTimeSendingBS = true
			sleep()
			return false, nil
		}
		if keyVal == IBUS_BackSpace {
			if e.nFakeBackSpace > 0 {
				e.nFakeBackSpace--
				return false, nil
			} else {
				e.preeditor.RemoveLastChar()
			}
			sleep()
			return false, nil
		}
	}
	// if the main thread is busy processing, the keypress events come all mixed up
	// so we enqueue these keypress events and process them sequentially on another thread
	keyPressChan <- [3]uint32{keyVal, keyCode, state}
	return true, nil
}

func (e *IBusBambooEngine) keyPressHandler(keyVal, keyCode, state uint32) {
	defer e.updateLastKeyWithShift(keyVal, state)
	if !e.isValidState(state) {
		e.preeditor.Reset()
		e.firstTimeSendingBS = true
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return
	}
	var keyRune = rune(keyVal)
	if keyVal == IBUS_BackSpace {
		if e.config.IBflags&IBautoNonVnRestore == 0 {
			if e.getRawKeyLen() > 0 {
				e.preeditor.RemoveLastChar()
			}
			e.ForwardKeyEvent(keyVal, keyCode, state)
			return
		}
		if e.getRawKeyLen() > 0 {
			oldRunes := []rune(e.getPreeditString())
			e.preeditor.RemoveLastChar()
			newRunes := []rune(e.getPreeditString())
			if len(oldRunes) == 0 {
				e.ForwardKeyEvent(keyVal, keyCode, state)
				return
			}
			e.updatePreviousText(newRunes, oldRunes, state)
			return
		}
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return
	}

	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		oldRunes := []rune(e.getPreeditString())
		e.preeditor.ProcessKey(keyRune, e.getMode())
		newRunes := []rune(e.getPreeditString())
		e.updatePreviousText(newRunes, oldRunes, state)
		return
	} else if bamboo.IsWordBreakSymbol(keyRune) {
		if keyVal == IBUS_Space && state&IBUS_SHIFT_MASK != 0 &&
			e.config.IBflags&IBrestoreKeyStrokesEnabled != 0 && !e.lastKeyWithShift {
			// restore key strokes
			var vnSeq = e.getPreeditString()
			if bamboo.HasVietnameseChar(vnSeq) {
				e.preeditor.RestoreLastWord()
				newRunes := []rune(e.getPreeditString())
				e.updatePreviousText(newRunes, []rune(vnSeq), state)
				return
			} else {
				e.SendText([]rune{keyRune})
				e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			}
			return
		}
		var processedStr = e.preeditor.GetProcessedString(bamboo.VietnameseMode)
		if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(processedStr) {
			// macro processing
			macText := e.expandMacro(processedStr)
			macText = macText + string(keyRune)
			e.updatePreviousText([]rune(macText), []rune(processedStr), state)
			e.preeditor.Reset()
			return
		}
		if e.mustFallbackToEnglish() {
			oldRunes := []rune(e.getPreeditString())
			e.preeditor.RestoreLastWord()
			newRunes := []rune(e.getComposedString())
			newRunes = append(newRunes, keyRune)
			e.updatePreviousText(newRunes, oldRunes, state)
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return
		}
		e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		e.SendText([]rune{keyRune})
		return
	}
	e.lastKeyWithShift = false
	e.preeditor.Reset()
	e.ForwardKeyEvent(keyVal, keyCode, state)
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
	log.Println("Updating Previous Text", string(oldRunes), string(newRunes), diffFrom)

	nBackSpace := 0
	// workaround for chrome and firefox's address bar
	if e.firstTimeSendingBS && diffFrom < newLen && diffFrom < oldLen && e.inBrowserList() && !e.inX11ShiftLeftList() {
		fmt.Println("Append a deadkey")
		e.SendText([]rune(" "))
		nBackSpace += 1
		time.Sleep(10 * time.Millisecond)
		e.firstTimeSendingBS = false
	}

	if diffFrom < oldLen {
		nBackSpace += oldLen - diffFrom
	}

	e.sendBackspaceAndNewRunes(nBackSpace, newRunes[diffFrom:])
}

func (e *IBusBambooEngine) sendBackspaceAndNewRunes(nBackSpace int, newRunes []rune) {
	if nBackSpace > 0 {
		if e.inXTestFakeKeyEventList() {
			e.nFakeBackSpace = nBackSpace
		} else if e.inX11ShiftLeftList() {
			e.nFakeShiftLeft = nBackSpace
		}
		e.SendBackSpace(nBackSpace)
	}
	e.SendText(newRunes)
}

func (e *IBusBambooEngine) SendBackSpace(n int) {
	// Gtk/Qt apps have a serious sync issue with fake backspaces
	// and normal string committing, so we'll not commit right now
	// but delay until all the sent backspaces got processed.
	if e.inXTestFakeKeyEventList() {
		var sleep = func() {
			var count = 0
			for e.nFakeBackSpace > 0 && count < 5 {
				time.Sleep(5 * time.Millisecond)
				count++
			}
			time.Sleep(20 * time.Millisecond)
		}
		fmt.Printf("Sendding %d backspace via XTestFakeKeyEvent\n", n)
		time.Sleep(20 * time.Millisecond)
		x11SendBackspace(n, 0)
		sleep()
	} else if e.inX11ShiftLeftList() {
		var sleep = func() {
			var count = 0
			for e.nFakeShiftLeft > 0 && count < 5 {
				time.Sleep(5 * time.Millisecond)
				count++
			}
		}
		fmt.Printf("Sendding %d Shift+Left via XTestFakeKeyEvent\n", n)
		time.Sleep(20 * time.Millisecond)
		x11SendShiftLeft(n, e.shiftRightIsPressing, 0)
		time.Sleep(time.Duration(n) * 20 * time.Millisecond)
		sleep()
	} else if e.inSurroundingTextList() {
		fmt.Printf("Sendding %d backspace via SurroundingText\n", n)
		e.DeleteSurroundingText(-int32(n), uint32(n))
		time.Sleep(20 * time.Millisecond)
	} else if e.inDirectForwardKeyList() {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("Sendding %d backspace via D_ForwardKeyEvent\n", n)
		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, IBUS_RELEASE_MASK)
			time.Sleep(5 * time.Millisecond)
		}
		time.Sleep(10 * time.Millisecond)
	} else if e.inForwardKeyList() {
		time.Sleep(10 * time.Millisecond)
		log.Printf("Sendding %d backspace via ForwardKeyEvent\n", n)

		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, IBUS_RELEASE_MASK)
			time.Sleep(5 * time.Millisecond)
		}
		time.Sleep(10 * time.Millisecond)
	} else {
		fmt.Println("There's something wrong with wmClasses")
	}
}

func (e *IBusBambooEngine) resetFakeBackspace() {
	e.nFakeBackSpace = 0
	e.nFakeShiftLeft = 0
}

func (e *IBusBambooEngine) SendText(rs []rune) {
	if len(rs) == 0 {
		return
	}
	if e.inDirectForwardKeyList() {
		log.Println("Forward as commit", string(rs))
		for _, chr := range rs {
			var keyVal = vnSymMapping[chr]
			if keyVal == 0 {
				keyVal = uint32(chr)
			}
			e.ForwardKeyEvent(keyVal, 0, 0)
			e.ForwardKeyEvent(keyVal, 0, IBUS_RELEASE_MASK)
			time.Sleep(5 * time.Millisecond)
		}
		return
	}
	log.Println("Sending text", string(rs))
	e.commitText(string(rs))
}
