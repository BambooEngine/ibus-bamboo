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
	"strings"
	"time"
	"unicode/utf8"

	"github.com/BambooEngine/bamboo-core"
	"github.com/godbus/dbus"
)

const BACKSPACE_INTERVAL = 0

func (e *IBusBambooEngine) bsProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var sleep = func() {
		log.Print("bsProcessKeyEvent sleeping....")
		for isProcessing || len(keyPressChan) > 0 {
			time.Sleep(5 * time.Millisecond)
		}
		log.Print("bsProcessKeyEvent awoke!")
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
			if e.nFakeBackSpace > 0 {
				e.nFakeBackSpace--
				return false, nil
			} else {
				sleep()
				if e.getRawKeyLen() > 0 {
					if e.shouldFallbackToEnglish(true) {
						e.preeditor.RestoreLastWord()
					}
					e.preeditor.RemoveLastChar(false)
				}
			}
			return false, nil
		}
		if keyVal == IBusTab {
			sleep()
			if ok, _ := e.getMacroText(); !ok {
				e.preeditor.Reset()
				return false, nil
			}
		}
	}
	if len(keyPressChan) == 0 && e.getRawKeyLen() == 0 && !inKeyList(e.preeditor.GetInputMethod().AppendingKeys, keyRune) {
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
	if !e.encounterControlKey {
		if e.isValidState(state) && e.canProcessKey(keyVal) {
			e.printableKeyCounter++
		} else {
			e.encounterControlKey = true
		}
	}
	// if the main thread is busy processing, the keypress events come all mixed up
	// so we enqueue these keypress events and process them sequentially on another thread
	keyPressChan <- [3]uint32{keyVal, keyCode, state}
	return true, nil
}

func (e *IBusBambooEngine) keyPressHandler(keyVal, keyCode, state uint32) {
	log.Printf(">>Backspace:ProcessKeyEvent >  %c | keyCode 0x%04x keyVal 0x%04x | %d\n", rune(keyVal), keyCode, keyVal, len(keyPressChan))
	defer e.updateLastKeyWithShift(keyVal, state)
	if e.keyPressDelay > 0 {
		time.Sleep(time.Duration(e.keyPressDelay) * time.Millisecond)
		e.keyPressDelay = 0
	}
	if !e.isValidState(state) || !e.canProcessKey(keyVal) {
		e.encounterControlKey = false
		e.printableKeyCounter = 0
	}
	if !e.isValidState(state) {
		e.preeditor.Reset()
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return
	}
	oldText := e.getPreeditString()
	_, oldMacText := e.getMacroText()
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
		if oldMacText != "" {
			e.updatePreviousText(oldMacText, oldText)
		} else {
			e.ForwardKeyEvent(keyVal, keyCode, state)
		}
		e.preeditor.Reset()
		return
	}

	newText, _ := e.getCommitText(keyVal, keyCode, state)
	if newText != "" {
		if e.shouldAppendDeadKey(newText, oldText) {
			fmt.Println("Append a deadkey")
			e.SendText([]rune(" "))
			time.Sleep(10 * time.Millisecond)
			e.isFirstTimeSendingBS = false
			e.SendBackSpace(1)
		}
		e.updatePreviousText(newText, oldText)
		return
	}
	e.preeditor.Reset()
	e.ForwardKeyEvent(keyVal, keyCode, state)
}

func (e *IBusBambooEngine) getCommitText(keyVal, keyCode, state uint32) (string, bool) {
	var keyRune = rune(keyVal)
	oldText := e.getPreeditString()
	_, oldMacText := e.getMacroText()
	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBusLockMask != 0 {
			keyRune = e.toUpper(keyRune)
		}
		e.preeditor.ProcessKey(keyRune, e.getBambooInputMode())
		if inKeyList(e.preeditor.GetInputMethod().AppendingKeys, keyRune) {
			if fullSeq := e.preeditor.GetProcessedString(bamboo.VietnameseMode); len(fullSeq) > 0 && rune(fullSeq[len(fullSeq)-1]) == keyRune {
				// u] => uo?
				return fullSeq, false
			} else if newText := e.getPreeditString(); newText != "" && keyRune == rune(newText[len(newText)-1]) {
				// ]] => ]
				e.preeditor.Reset()
				return oldText + string(keyRune), true
			} else {
				// ] => o?
				return e.getPreeditString(), false
			}
		} else {
			return e.getPreeditString(), false
		}
	} else if bamboo.IsWordBreakSymbol(keyRune) {
		// restore key strokes by pressing Shift + Space
		if keyVal == IBusSpace && state&IBusShiftMask != 0 &&
			e.config.IBflags&IBrestoreKeyStrokesEnabled != 0 && !e.lastKeyWithShift {
			if bamboo.HasAnyVietnameseRune(oldText) {
				commitText := e.preeditor.GetProcessedString(bamboo.EnglishMode)
				e.preeditor.RestoreLastWord()
				return commitText, false
			}
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return oldText + string(keyRune), true
		}
		// macro processing
		if oldMacText != "" {
			macText := oldMacText + string(keyRune)
			e.preeditor.Reset()
			return macText, true
		}
		if bamboo.HasAnyVietnameseRune(oldText) && e.mustFallbackToEnglish() {
			e.preeditor.RestoreLastWord()
			newText := e.preeditor.GetProcessedString(bamboo.EnglishMode) + string(keyRune)
			e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
			return newText, true
		}
		e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
		fmt.Println("cannot process key rune is ", keyRune)
		log.Printf("englis [%s]\n", e.preeditor.GetProcessedString(bamboo.EnglishMode))
		return oldText + string(keyRune), true
	}
	return "", false
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

func (e *IBusBambooEngine) shouldAppendDeadKey(newText, oldText string) bool {
	var oldRunes = []rune(oldText)
	var newRunes = []rune(newText)
	var offset = e.getPreeditOffset(newRunes, oldRunes)

	// workaround for chrome and firefox's address bar
	if e.isFirstTimeSendingBS && offset < len(newRunes) && offset < len(oldRunes) && e.inBrowserList() &&
		!e.checkInputMode(shiftLeftForwardingIM) {
		return true
	}
	return false
}

func (e *IBusBambooEngine) updatePreviousText(newText, oldText string) {
	offsetRunes, nBackSpace := e.getOffsetRunes(newText, oldText)
	e.printableKeyCounter--
	if nBackSpace > 0 {
		e.SendBackSpace(nBackSpace)
		var history = []string{string(offsetRunes)}
		var i int
		for i = 0; i < e.printableKeyCounter && i < len(keyPressChan); i++ {
			var data = <-keyPressChan
			var commitText, isLast = e.getCommitText(data[0], data[1], data[2])
			history[len(history)-1] = commitText
			if isLast {
				history = append(history, "")
			}
		}
		if i > 0 {
			e.printableKeyCounter -= i
			fmt.Print("\n\nHISTORY\n====================================\n         ", history, len(history))
			fmt.Print("\n====================================\n\n")
			fullRunes := []rune(strings.Join(history, ""))
			offsetRunes0, nBackSpace0 := e.getOffsetRunes(strings.Join(history, ""), oldText)
			if nBackSpace0 > nBackSpace {
				e.SendBackSpace(nBackSpace0 - nBackSpace)
			} else if nBackSpace0 < nBackSpace {
				var offset = utf8.RuneCountInString(oldText) - nBackSpace
				offsetRunes0 = fullRunes[offset:]
			}
			log.Printf("Updating Previous Text %s ---> %s\n", oldText, string(fullRunes))
			e.SendText(offsetRunes0)
			return
		}
	}
	log.Printf("Updating Previous Text %s ---> %s\n", oldText, newText)
	e.SendText(offsetRunes)
}

func (e *IBusBambooEngine) getOffsetRunes(newText, oldText string) ([]rune, int) {
	var oldRunes = []rune(oldText)
	var newRunes = []rune(newText)
	var nBackSpace = 0
	var offset = e.getPreeditOffset(newRunes, oldRunes)
	if offset < len(oldRunes) {
		nBackSpace += len(oldRunes) - offset
	}

	return newRunes[offset:], nBackSpace
}

func (e *IBusBambooEngine) SendBackSpace(n int) {
	// Gtk/Qt apps have a serious sync issue with fake backspaces
	// and normal string committing, so we'll not commit right now
	// but delay until all the sent backspaces got processed.
	if e.checkInputMode(xTestFakeKeyEventIM) {
		e.nFakeBackSpace = n
		var sleep = func() {
			var count = 0
			for e.nFakeBackSpace > 0 && count < 10 {
				time.Sleep(5 * time.Millisecond)
				count++
			}
		}
		fmt.Printf("Sendding %d backspace via XTestFakeKeyEvent\n", n)
		time.Sleep(30 * time.Millisecond)
		x11SendBackspace(n, 0)
		sleep()
		time.Sleep(time.Duration(n) * (30 + BACKSPACE_INTERVAL) * time.Millisecond)
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
		time.Sleep(time.Duration(n) * (20 + BACKSPACE_INTERVAL) * time.Millisecond)
	} else if e.checkInputMode(shiftLeftForwardingIM) {
		time.Sleep(30 * time.Millisecond)
		log.Printf("Sendding %d Shift+Left via shiftLeftForwardingIM\n", n)

		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBusLeft, XkLeft-8, IBusShiftMask)
			e.ForwardKeyEvent(IBusLeft, XkLeft-8, IBusReleaseMask)
		}
		time.Sleep(time.Duration(n) * (30 + BACKSPACE_INTERVAL) * time.Millisecond)
	} else if e.checkInputMode(backspaceForwardingIM) {
		time.Sleep(30 * time.Millisecond)
		log.Printf("Sendding %d backspace via backspaceForwardingIM\n", n)

		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBusBackSpace, XkBackspace-8, 0)
			e.ForwardKeyEvent(IBusBackSpace, XkBackspace-8, IBusReleaseMask)
		}
		time.Sleep(time.Duration(n) * (30 + BACKSPACE_INTERVAL) * time.Millisecond)
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
