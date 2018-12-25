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
	"strings"
	"time"
)

func (e *IBusBambooEngine) preeditProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var rawKeyLen = e.getRawKeyLen()

	if keyVal == IBUS_BackSpace {
		e.ignorePreedit = false
		if rawKeyLen > 0 {
			e.preediter.RemoveLastChar()
			e.updatePreedit()
			return true, nil
		} else {
			return false, nil
		}
	}

	if keyVal == IBUS_space || keyVal == IBUS_KP_Space {
		e.ignorePreedit = false
		var processedStr = e.getComposedString()
		if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(processedStr) {
			processedStr = e.macroTable.GetText(processedStr)
			e.preediter.Reset()
			e.commitText(e.encodeText(processedStr))
		} else {
			e.commitPreedit(0)
		}
		return false, nil
	}

	if keyVal == IBUS_Return || keyVal == IBUS_KP_Enter {
		e.ignorePreedit = false
		if rawKeyLen > 0 {
			e.commitPreedit(keyVal)
			e.ForwardKeyEvent(keyVal, keyCode, state)
			return true, nil
		} else {
			return false, nil
		}
	}

	if keyVal == IBUS_Escape {
		e.ignorePreedit = false
		if rawKeyLen > 0 {
			e.commitPreedit(keyVal)
			return true, nil
		}
		return false, nil
	}

	if (keyVal >= 'a' && keyVal <= 'z') ||
		(keyVal >= 'A' && keyVal <= 'Z') ||
		(keyVal >= '0' && keyVal <= '9') ||
		(inKeyMap(e.preediter.GetInputMethod().Keys, rune(keyVal))) {
		var keyRune = rune(keyVal)
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		if e.ignorePreedit {
			return false, nil
		}
		if e.config.IBflags&IBautoNonVnRestore == 0 {
			e.preediter.ProcessChar(keyRune, bamboo.VietnameseMode)
			e.updatePreedit()
			return true, nil
		}
		e.preediter.ProcessChar(keyRune, e.getMode())
		if (e.config.IBflags&IBautoCommitWithVnNotMatch != 0 &&
			e.getSpellingMatchResult(false) == bamboo.FindResultNotMatch) ||
			(e.config.IBflags&IBautoCommitWithVnFullMatch != 0 && e.preediter.HasTone() &&
				e.getSpellingMatchResult(true) == bamboo.FindResultMatchFull) {
			e.ignorePreedit = true
			e.commitPreedit(0)
			e.preediter.Reset()
			return true, nil
		}
		e.updatePreedit()
		return true, nil
	} else {
		e.commitPreedit(keyVal)
		return false, nil
	}
	return false, nil
}

var preeditUpdateChan = make(chan uint32)

func (e *IBusBambooEngine) startAutoCommit() {
	for {
		var timeout = e.config.AutoCommitAfter
		select {
		case <-preeditUpdateChan:
			break
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			var rawKeyLen = e.getRawKeyLen()
			if e.config.IBflags&IBautoCommitWithDelay != 0 && rawKeyLen > 0 {
				e.commitPreedit(0)
			}
		}
	}
}

func (e *IBusBambooEngine) updatePreedit() {
	var processedStr = e.getPreeditString()
	var preeditLen = uint32(len([]rune(processedStr)))
	var ibusText = ibus.NewText(processedStr)

	if e.config.IBflags&IBpreeditInvisibility != 0 {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_NONE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	} else {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	}

	e.UpdatePreeditTextWithMode(ibusText, preeditLen, true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	if preeditLen == 0 {
		e.HidePreeditText()
		e.preediter.Reset()
	}
	mouseCaptureUnlock()

	preeditUpdateChan <- 0
}

func (e *IBusBambooEngine) shouldFallbackToEnglish() bool {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return false
	}
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return false
	}
	if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(vnSeq) {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && (vnRunes[len(vnRunes)-1] == 'd' || strings.ContainsRune(vnSeq, 'đ')) {
		if !bamboo.HasVowel(vnRunes) {
			return false
		}
	}
	if e.getSpellingMatchResult(false) != bamboo.FindResultNotMatch {
		return false
	}
	return true
}

func (e *IBusBambooEngine) mustFallbackToEnglish() bool {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return false
	}
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && strings.ContainsRune(vnSeq, 'đ') {
		if !bamboo.HasVowel(vnRunes) {
			return false
		}
	}
	if e.config.IBflags&IBspellCheckingWithDicts != 0 {
		return !e.dictionary[strings.ToLower(vnSeq)]
	}
	if e.getSpellingMatchResult(false) == bamboo.FindResultMatchFull {
		return false
	}
	return true
}

func (e *IBusBambooEngine) isSpellingCorrect() bool {
	return e.getSpellingMatchResult(false) == bamboo.FindResultMatchFull
}

func (e *IBusBambooEngine) getSpellingMatchResult(deepSearch bool) uint8 {
	return e.preediter.GetSpellingMatchResult(bamboo.NoTone, deepSearch)
}

func (e *IBusBambooEngine) getComposedString() string {
	var processedStr string
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		processedStr = e.getVnSeq()
	} else {
		processedStr = e.getProcessedString(bamboo.VietnameseMode)
	}
	if e.mustFallbackToEnglish() {
		processedStr = e.getProcessedString(bamboo.EnglishMode)
		return processedStr
	}
	return processedStr
}

func (e *IBusBambooEngine) encodeText(text string) string {
	return bamboo.Encode(e.config.Charset, text)
}

func (e *IBusBambooEngine) getProcessedString(mode bamboo.Mode) string {
	if e.config.IBflags&IBautoNonVnRestore != 0 {
		return e.preediter.GetProcessedString(mode)
	}
	return e.getVnSeq()
}

func (e *IBusBambooEngine) getPreeditString() string {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return e.getVnSeq()
	}
	if e.shouldFallbackToEnglish() {
		return e.getProcessedString(bamboo.EnglishMode)
	}
	return e.getProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) getMode() bamboo.Mode {
	if e.shouldFallbackToEnglish() {
		return bamboo.EnglishMode
	}
	return bamboo.VietnameseMode
}

func (e *IBusBambooEngine) commitPreedit(lastKey uint32) {
	var commitStr string
	commitStr += e.getComposedString()

	e.preediter.Reset()
	e.commitText(e.encodeText(commitStr))
}

func (e *IBusBambooEngine) commitText(str string) {
	for _, chr := range []rune(str) {
		e.CommitText(ibus.NewText(string(chr)))
	}
	//e.CommitText(ibus.NewText(commitStr))

	e.HidePreeditText()
}

func (e *IBusBambooEngine) getVnSeq() string {
	return e.preediter.GetProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) hasMacroKey(key string) bool {
	return e.macroTable.GetText(key) != ""
}
