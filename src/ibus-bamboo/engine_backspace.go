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

var keyPressChan = make(chan [3]uint32, 100)

func (e *IBusBambooEngine) bsProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	if e.inXTestFakeKeyEventList() || e.inSurroundingTextList() {
		// we don't want to use ForwardKeyEvent Api in X11 XTestFakeKeyEvent and Surrounding Text mode
		var sleep = func() {
			for len(keyPressChan) > 0 {
				time.Sleep(5 * time.Millisecond)
			}
		}
		if !e.isValidState(state) || !e.canProcessKey(keyVal) {
			e.preeditor.Reset()
			e.resetFakeBackspace()
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

func (e *IBusBambooEngine) keyPressHandler() {
	for {
		select {
		case keyEvents := <-keyPressChan:
			var keyVal, keyCode, state = keyEvents[0], keyEvents[1], keyEvents[2]
			if !e.isValidState(state) {
				e.preeditor.Reset()
				e.ForwardKeyEvent(keyVal, keyCode, state)
				break
			}
			var keyRune = rune(keyVal)
			if keyVal == IBUS_BackSpace {
				if e.config.IBflags&IBautoNonVnRestore == 0 {
					if e.getRawKeyLen() > 0 {
						e.preeditor.RemoveLastChar()
					}
					e.ForwardKeyEvent(keyVal, keyCode, state)
					break
				}
				if e.getRawKeyLen() > 0 {
					oldRunes := []rune(e.getPreeditString())
					e.preeditor.RemoveLastChar()
					newRunes := []rune(e.getPreeditString())
					if len(oldRunes) == 0 {
						e.ForwardKeyEvent(keyVal, keyCode, state)
						break
					}
					e.updatePreviousText(newRunes, oldRunes, state)
					break
				}
				e.ForwardKeyEvent(keyVal, keyCode, state)
				break
			}

			if e.preeditor.CanProcessKey(keyRune) {
				if state&IBUS_LOCK_MASK != 0 {
					keyRune = toUpper(keyRune)
				}
				oldRunes := []rune(e.getPreeditString())
				e.preeditor.ProcessKey(keyRune, e.getMode())
				newRunes := []rune(e.getPreeditString())
				e.updatePreviousText(newRunes, oldRunes, state)
				break
			} else if bamboo.IsWordBreakSymbol(keyRune) {
				// macro processing
				var processedStr = e.preeditor.GetProcessedString(bamboo.VietnameseMode, true)
				if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(processedStr) {
					macText := e.macroTable.GetText(processedStr)
					macText = macText + string(keyRune)
					e.updatePreviousText([]rune(macText), []rune(processedStr), state)
					e.preeditor.Reset()
					break
				} else if e.mustFallbackToEnglish() && !e.inXTestFakeKeyEventList() && !e.inSurroundingTextList() {
					oldRunes := []rune(e.getPreeditString())
					newRunes := []rune(e.getComposedString())
					newRunes = append(newRunes, keyRune)
					e.updatePreviousText(newRunes, oldRunes, state)
					e.preeditor.Reset()
					break
				}
				e.preeditor.ProcessKey(keyRune, bamboo.EnglishMode)
				e.SendText([]rune{keyRune})
				break
			}
			e.preeditor.Reset()
			e.ForwardKeyEvent(keyVal, keyCode, state)
		}
	}
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
	fmt.Println("Updating Previous Text", string(oldRunes), string(newRunes), diffFrom)

	nBackSpace := 0
	// workaround for chrome and firefox's address bar
	if e.firstTimeSendingBS && diffFrom < newLen && diffFrom < oldLen && e.inBrowserList() {
		fmt.Println("Append a deadkey")
		e.SendText([]rune(" "))
		time.Sleep(10 * time.Millisecond)
		nBackSpace += 1
		e.firstTimeSendingBS = false
	}

	if diffFrom < oldLen {
		nBackSpace += oldLen - diffFrom
	}

	e.sendBackspaceAndNewRunes(nBackSpace, newRunes[diffFrom:])
	mouseCaptureUnlock()
}

func (e *IBusBambooEngine) sendBackspaceAndNewRunes(nBackSpace int, newRunes []rune) {
	if e.inXTestFakeKeyEventList() {
		var sleep = func() {
			var count = 0
			for e.nFakeBackSpace > 0 && count < 5 {
				time.Sleep(5 * time.Millisecond)
				count++
			}
			time.Sleep(20 * time.Millisecond)
		}
		if nBackSpace > 0 {
			e.nFakeBackSpace = nBackSpace
			e.SendBackSpace(nBackSpace)
			sleep()
			e.SendText(newRunes)
		} else {
			e.SendText(newRunes)
		}
		return
	}
	if nBackSpace > 0 {
		e.SendBackSpace(nBackSpace)
	}
	e.SendText(newRunes)
}

func (e *IBusBambooEngine) SendBackSpace(n int) {
	if e.inXTestFakeKeyEventList() {
		time.Sleep(20 * time.Millisecond)
		fmt.Printf("Sendding %d backspace via XTestFakeKeyEvent\n", n)
		if e.inChromeList() { // workaround for chrome's address bar
			x11SendBackspace(n, 0)
			time.Sleep(time.Duration(n) * 10 * time.Millisecond)
		} else {
			x11SendBackspace(n, 10)
		}
	} else if e.inSurroundingTextList() {
		time.Sleep(20 * time.Millisecond)
		fmt.Printf("Sendding %d backspace via SurroundingText\n", n)
		e.DeleteSurroundingText(-int32(n), uint32(n))
		time.Sleep(20 * time.Millisecond)
	} else if e.inDirectForwardKeyList() {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("Sendding %d backspace via ForwardKeyEvent *\n", n)
		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, 0x16-8, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, 0x16-8, IBUS_RELEASE_MASK)
			time.Sleep(5 * time.Millisecond)
		}
		time.Sleep(10 * time.Millisecond)
	} else if e.inForwardKeyList() {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("Sendding %d backspace via ForwardKeyEvent **\n", n)
		for i := 0; i < n; i++ {
			e.ForwardKeyEvent(IBUS_BackSpace, 0x16-8, 0)
			e.ForwardKeyEvent(IBUS_BackSpace, 0x16-8, IBUS_RELEASE_MASK)
			time.Sleep(5 * time.Millisecond)
		}
		time.Sleep(10 * time.Millisecond)
	} else {
		fmt.Println("There's something wrong with wmClasses")
	}
}

func (e *IBusBambooEngine) resetFakeBackspace() {
	e.nFakeBackSpace = 0
}

func (e *IBusBambooEngine) SendText(rs []rune) {
	if e.inDirectForwardKeyList() {
		log.Println("Forward as commit", string(rs))
		for _, chr := range rs {
			var keyVal = keysymsMapping[chr]
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
	e.CommitText(ibus.NewText(string(rs)))
}
