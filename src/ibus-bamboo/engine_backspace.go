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
	if isMovementKey(keyVal) {
		e.preeditor.Reset()
		e.resetFakeBackspace()
		e.isSurroundingTextReady = true
		return false, nil
	}
	if e.checkInputMode(xTestFakeKeyEventIM) || e.checkInputMode(surroundingTextIM) {
		// we don't want to use ForwardKeyEvent api in X11 XTestFakeKeyEvent and Surrounding Text mode
		var sleep = func() {
			for len(keyPressChan) > 0 {
				time.Sleep(5 * time.Millisecond)
			}
		}
		if keyVal == IBUS_Left && state&IBUS_SHIFT_MASK != 0 {
			return false, nil
		}
		if !e.isValidState(state) || !e.canProcessKey(keyVal, state) {
			e.preeditor.Reset()
			e.resetFakeBackspace()
			e.isFirstTimeSendingBS = true
			sleep()
			return false, nil
		}
		if keyVal == IBUS_BackSpace {
			if e.nFakeBackSpace > 0 {
				e.nFakeBackSpace--
				return false, nil
			} else {
				e.preeditor.RemoveLastChar(false)
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
	if e.keyPressDelay > 0 {
		time.Sleep(time.Duration(e.keyPressDelay) * time.Millisecond)
		e.keyPressDelay = 0
	}
	if !e.isValidState(state) {
		e.preeditor.Reset()
		e.isFirstTimeSendingBS = true
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return
	}
	var keyRune = rune(keyVal)
	oldText := e.preeditor.GetProcessedString(bamboo.VietnameseMode | bamboo.WithEffectKeys)
	if keyVal == IBUS_BackSpace {
		if e.config.IBflags&IBautoNonVnRestore == 0 || e.checkInputMode(shiftLeftForwardingIM) {
			if e.getRawKeyLen() > 0 {
				e.preeditor.RemoveLastChar(false)
			}
			e.ForwardKeyEvent(keyVal, keyCode, state)
			return
		}
		if e.getRawKeyLen() > 0 {
			e.preeditor.RemoveLastChar(true)
			if oldText == "" {
				e.ForwardKeyEvent(keyVal, keyCode, state)
				return
			}
			e.updatePreviousText(e.preeditor.GetProcessedString(bamboo.VietnameseMode|bamboo.WithEffectKeys), oldText)
			return
		}
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return
	}

	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		e.preeditor.ProcessKey(keyRune, e.getInputMethod())
		var vnSeq = e.preeditor.GetProcessedString(bamboo.VietnameseMode | bamboo.WithEffectKeys)
		if len(vnSeq) > 0 && rune(vnSeq[len(vnSeq)-1]) == keyRune && bamboo.IsWordBreakSymbol(keyRune) {
			e.updatePreviousText(vnSeq, oldText)
			e.preeditor.Reset()
		} else {
			e.updatePreviousText(vnSeq, oldText)
		}
		return
	} else if bamboo.IsWordBreakSymbol(keyRune) || ('0' <= keyVal && keyVal <= '9') {
		if keyVal == IBUS_Space && state&IBUS_SHIFT_MASK != 0 &&
			e.config.IBflags&IBrestoreKeyStrokesEnabled != 0 && !e.lastKeyWithShift {
			// restore key strokes
			if bamboo.HasAnyVietnameseRune(oldText) {
				e.preeditor.RestoreLastWord()
				e.updatePreviousText(e.getPreeditString(), oldText)
				return
			} else {
				e.SendText([]rune{keyRune})
				e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			}
			return
		}
		if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(oldText) {
			// macro processing
			macText := e.expandMacro(oldText)
			macText = macText + string(keyRune)
			e.updatePreviousText(macText, oldText)
			e.preeditor.Reset()
			return
		}
		if e.mustFallbackToEnglish() {
			e.preeditor.RestoreLastWord()
			newText := e.preeditor.GetProcessedString(bamboo.EnglishMode|bamboo.WithEffectKeys) + string(keyRune)
			e.updatePreviousText(newText, oldText)
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return
		}
		e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		e.SendText([]rune{keyRune})
		return
	}
	e.preeditor.Reset()
	e.ForwardKeyEvent(keyVal, keyCode, state)
}

func (e *IBusBambooEngine) updatePreviousText(newText, oldText string) {
	var oldRunes = []rune(oldText)
	var newRunes = []rune(newText)
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
	log.Printf("Updating Previous Text %s ---> %s\n", string(oldRunes), string(newRunes))

	nBackSpace := 0
	// workaround for chrome and firefox's address bar
	if e.isFirstTimeSendingBS && diffFrom < newLen && diffFrom < oldLen && e.inBrowserList() &&
		!e.checkInputMode(shiftLeftForwardingIM) {
		fmt.Println("Append a deadkey")
		e.SendText([]rune(" "))
		nBackSpace += 1
		time.Sleep(10 * time.Millisecond)
		e.isFirstTimeSendingBS = false
	}

	if diffFrom < oldLen {
		nBackSpace += oldLen - diffFrom
	}

	e.sendBackspaceAndNewRunes(nBackSpace, newRunes[diffFrom:])
}

func (e *IBusBambooEngine) sendBackspaceAndNewRunes(nBackSpace int, newRunes []rune) {
	if nBackSpace > 0 {
		if e.checkInputMode(xTestFakeKeyEventIM) {
			e.nFakeBackSpace = nBackSpace
		}
		e.SendBackSpace(nBackSpace)
	}
	e.SendText(newRunes)
}

func (e *IBusBambooEngine) SendBackSpace(n int) {
	// Gtk/Qt apps have a serious sync issue with fake backspaces
	// and normal string committing, so we'll not commit right now
	// but delay until all the sent backspaces got processed.
	if e.checkInputMode(xTestFakeKeyEventIM) {
		var sleep = func() {
			var count = 0
			for e.nFakeBackSpace > 0 && count < 5 {
				time.Sleep(5 * time.Millisecond)
				count++
			}
		}
		fmt.Printf("Sendding %d backspace via XTestFakeKeyEvent\n", n)
		time.Sleep(30 * time.Millisecond)
		x11SendBackspace(n, 0)
		sleep()
		time.Sleep(time.Duration(n) * 30 * time.Millisecond)
	} else if e.checkInputMode(surroundingTextIM) {
		time.Sleep(20 * time.Millisecond)
		fmt.Printf("Sendding %d backspace via SurroundingText\n", n)
		e.DeleteSurroundingText(-int32(n), uint32(n))
		time.Sleep(20 * time.Millisecond)
	} else if e.checkInputMode(forwardAsCommitIM) {
		time.Sleep(20 * time.Millisecond)
		fmt.Printf("Sendding %d backspace via forwardAsCommitIM\n", n)
		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, IBUS_RELEASE_MASK)
		}
		time.Sleep(time.Duration(n) * 30 * time.Millisecond)
	} else if e.checkInputMode(shiftLeftForwardingIM) {
		time.Sleep(30 * time.Millisecond)
		log.Printf("Sendding %d Shift+Left via shiftLeftForwardingIM\n", n)

		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_Left, XK_Left-8, IBUS_SHIFT_MASK)
			e.ForwardKeyEvent(IBUS_Left, XK_Left-8, IBUS_RELEASE_MASK)
		}
		time.Sleep(time.Duration(n) * 30 * time.Millisecond)
	} else if e.checkInputMode(backspaceForwardingIM) {
		time.Sleep(30 * time.Millisecond)
		log.Printf("Sendding %d backspace via backspaceForwardingIM\n", n)

		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, XK_BackSpace-8, IBUS_RELEASE_MASK)
		}
		time.Sleep(time.Duration(n) * 30 * time.Millisecond)
	} else {
		fmt.Println("There's something wrong with wmClasses")
	}
}

func (e *IBusBambooEngine) resetFakeBackspace() {
	e.nFakeBackSpace = 0
}

func (e *IBusBambooEngine) SendText(rs []rune) {
	if len(rs) == 0 {
		return
	}
	if e.checkInputMode(forwardAsCommitIM) {
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
	e.commitText(string(rs))
}
