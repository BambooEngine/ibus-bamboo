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
	"log"
	"time"

	"github.com/BambooEngine/bamboo-core"
	"github.com/godbus/dbus"
)

func (e *IBusBambooEngine) bsProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var sleep = func() {
		for len(keyPressChan) > 0 {
			time.Sleep(5 * time.Millisecond)
		}
	}
	if isMovementKey(keyVal) {
		e.preeditor.Reset()
		e.resetFakeBackspace()
		e.isSurroundingTextReady = true
		return false, nil
	}
	var keyRune = rune(keyVal)
	// Caution: don't use ForwardKeyEvent api in XTestFakeKeyEvent and SurroundingText mode
	if e.checkInputMode(xTestFakeKeyEventIM) || e.checkInputMode(surroundingTextIM) {
		if keyVal == IBusLeft && state&IBusShiftMask != 0 {
			return false, nil
		}
		if !e.isValidState(state) || !e.canProcessKey(keyVal) {
			sleep()
			e.preeditor.Reset()
			e.resetFakeBackspace()
			return false, nil
		}
		if keyVal == IBusBackSpace {
			sleep()
			if e.nFakeBackSpace > 0 {
				e.nFakeBackSpace--
				return false, nil
			} else {
				e.preeditor.RemoveLastChar(false)
			}
			return false, nil
		}
		if keyVal == IBusTab {
			sleep()
			oldText := e.preeditor.GetProcessedString(bamboo.VietnameseMode)
			if e.config.IBflags&IBmacroEnabled != 0 && !e.macroTable.HasKey(oldText) {
				e.preeditor.Reset()
				return false, nil
			}
		}
	}
	if e.getRawKeyLen() == 0 && !inKeyList(e.preeditor.GetInputMethod().AppendingKeys, keyRune) {
		sleep()
		if e.getRawKeyLen() == 0 {
			e.updateLastKeyWithShift(keyVal, state)
			if e.preeditor.CanProcessKey(keyRune) && e.isValidState(state) {
				e.isFirstTimeSendingBS = true
				if state&IBusLockMask != 0 {
					keyRune = e.toUpper(keyRune)
				}
				e.preeditor.ProcessKey(keyRune, bamboo.VietnameseMode)
			}
			return false, nil
		}
	}
	// if the main thread is busy processing, the keypress events come all mixed up
	// so we enqueue these keypress events and process them sequentially on another thread
	keyPressChan <- [3]uint32{keyVal, keyCode, state}
	return true, nil
}

func (e *IBusBambooEngine) keyPressHandler(keyVal, keyCode, state uint32) {
	log.Printf("Backspace:ProcessKeyEvent >  %c | keyCode 0x%04x keyVal 0x%04x | %d\n", rune(keyVal), keyCode, keyVal, len(keyPressChan))
	defer e.updateLastKeyWithShift(keyVal, state)
	if e.keyPressDelay > 0 {
		time.Sleep(time.Duration(e.keyPressDelay) * time.Millisecond)
		e.keyPressDelay = 0
	}
	if !e.isValidState(state) {
		e.preeditor.Reset()
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return
	}
	var keyRune = rune(keyVal)
	oldText := e.getPreeditString()
	if keyVal == IBusBackSpace {
		if e.getRawKeyLen() > 0 {
			if e.config.IBflags&IBautoNonVnRestore == 0 {
				e.preeditor.RemoveLastChar(false)
				e.ForwardKeyEvent(keyVal, keyCode, state)
				return
			}
			e.preeditor.RemoveLastChar(true)
			var newText = e.getPreeditString()
			var offset = e.getPreeditOffset([]rune(newText), []rune(oldText))
			if oldText != "" && offset != len([]rune(newText)) {
				e.updatePreviousText(newText, oldText)
				return
			}
		}
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return
	}

	if keyVal == IBusTab {
		if e.config.IBflags&IBmacroEnabled != 0 && e.macroTable.HasKey(oldText) {
			// macro processing
			macText := e.expandMacro(oldText)
			e.updatePreviousText(macText, oldText)
		} else {
			e.ForwardKeyEvent(keyVal, keyCode, state)
		}
		e.preeditor.Reset()
		return
	}

	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBusLockMask != 0 {
			keyRune = e.toUpper(keyRune)
		}
		e.preeditor.ProcessKey(keyRune, e.getBambooInputMode())
		if inKeyList(e.preeditor.GetInputMethod().AppendingKeys, keyRune) {
			if fullSeq := e.preeditor.GetProcessedString(bamboo.VietnameseMode); len(fullSeq) > 0 && rune(fullSeq[len(fullSeq)-1]) == keyRune {
				e.updatePreviousText(fullSeq, oldText)
				e.preeditor.Reset()
			} else if newText := e.getPreeditString(); newText != "" && keyRune == rune(newText[len(newText)-1]) {
				e.SendText([]rune{keyRune})
				e.preeditor.Reset()
			} else {
				e.updatePreviousText(e.getPreeditString(), oldText)
			}
		} else {
			e.updatePreviousText(e.getPreeditString(), oldText)
		}
		return
	} else if bamboo.IsWordBreakSymbol(keyRune) {
		if keyVal == IBusSpace && state&IBusShiftMask != 0 &&
			e.config.IBflags&IBrestoreKeyStrokesEnabled != 0 && !e.lastKeyWithShift {
			// restore key strokes
			if bamboo.HasAnyVietnameseRune(oldText) {
				e.updatePreviousText(e.preeditor.GetProcessedString(bamboo.EnglishMode), oldText)
				e.preeditor.RestoreLastWord()
				return
			} else {
				e.SendText([]rune{keyRune})
				e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			}
			return
		}
		if e.config.IBflags&IBmacroEnabled != 0 && e.macroTable.HasKey(oldText) {
			// macro processing
			macText := e.expandMacro(oldText)
			macText = macText + string(keyRune)
			e.updatePreviousText(macText, oldText)
			e.preeditor.Reset()
			return
		}
		if bamboo.HasAnyVietnameseRune(oldText) && e.mustFallbackToEnglish() {
			e.preeditor.RestoreLastWord()
			newText := e.preeditor.GetProcessedString(bamboo.EnglishMode) + string(keyRune)
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

func (e *IBusBambooEngine) getPreeditOffset(newRunes, oldRunes []rune) int {
	var minLen = len(oldRunes)
	if len(newRunes) < minLen {
		minLen = len(newRunes)
	}
	for i := 0; i < minLen; i++ {
		if oldRunes[i] != newRunes[i] {
			return i
		}
	}
	return minLen
}

func (e *IBusBambooEngine) updatePreviousText(newText, oldText string) {
	var oldRunes = []rune(oldText)
	var newRunes = []rune(newText)
	var nBackSpace = 0
	var offset = e.getPreeditOffset(newRunes, oldRunes)
	if offset < len(oldRunes) {
		nBackSpace += len(oldRunes) - offset
	}

	// workaround for chrome and firefox's address bar
	if e.isFirstTimeSendingBS && offset < len(newRunes) && offset < len(oldRunes) && e.inBrowserList() &&
		!e.checkInputMode(shiftLeftForwardingIM) {
		fmt.Println("Append a deadkey")
		e.SendText([]rune(" "))
		nBackSpace += 1
		time.Sleep(10 * time.Millisecond)
		e.isFirstTimeSendingBS = false
	}

	log.Printf("Updating Previous Text %s ---> %s\n", oldText, newText)
	e.sendBackspaceAndNewRunes(nBackSpace, newRunes[offset:])
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
			e.ForwardKeyEvent(IBusBackSpace, XkBackspace-8, 0)
			e.ForwardKeyEvent(IBusBackSpace, XkBackspace-8, IBusReleaseMask)
		}
		time.Sleep(time.Duration(n) * 20 * time.Millisecond)
	} else if e.checkInputMode(shiftLeftForwardingIM) {
		time.Sleep(30 * time.Millisecond)
		log.Printf("Sendding %d Shift+Left via shiftLeftForwardingIM\n", n)

		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBusLeft, XkLeft-8, IBusShiftMask)
			e.ForwardKeyEvent(IBusLeft, XkLeft-8, IBusReleaseMask)
		}
		time.Sleep(time.Duration(n) * 30 * time.Millisecond)
	} else if e.checkInputMode(backspaceForwardingIM) {
		time.Sleep(30 * time.Millisecond)
		log.Printf("Sendding %d backspace via backspaceForwardingIM\n", n)

		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBusBackSpace, XkBackspace-8, 0)
			e.ForwardKeyEvent(IBusBackSpace, XkBackspace-8, IBusReleaseMask)
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
			e.ForwardKeyEvent(keyVal, 0, IBusReleaseMask)
		}
		time.Sleep(time.Duration(len(rs)) * 5 * time.Millisecond)
		return
	}
	e.commitText(string(rs))
}
