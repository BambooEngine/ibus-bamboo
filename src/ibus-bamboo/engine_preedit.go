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
	"strings"
)

func (e *IBusBambooEngine) preeditProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var rawKeyLen = e.getRawKeyLen()
	var keyRune = rune(keyVal)

	if !e.isValidState(state) {
		e.commitPreedit()
		return false, nil
	}
	if keyVal == IBUS_Caps_Lock ||
		(!(state&IBUS_SHIFT_MASK != 0) && (keyVal == IBUS_Shift_L || keyVal == IBUS_Shift_R)) { // when press one shift key
		return false, nil
	}

	if keyVal == IBUS_BackSpace {
		if rawKeyLen > 0 {
			e.preeditor.RemoveLastChar()
			e.updatePreedit()
			return true, nil
		} else {
			return false, nil
		}
	}
	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		if e.ignorePreedit {
			return false, nil
		}
		e.preeditor.ProcessKey(keyRune, e.getMode())
		if (e.config.IBflags&IBautoCommitWithVnNotMatch != 0 &&
			e.getSpellingMatchResult(false) == bamboo.FindResultNotMatch) ||
			(e.config.IBflags&IBautoCommitWithVnFullMatch != 0 && e.preeditor.HasTone() &&
				e.getSpellingMatchResult(true) == bamboo.FindResultMatchFull) {
			e.ignorePreedit = true
			e.commitPreedit()
			return true, nil
		}
		e.updatePreedit()
		return true, nil
	} else if bamboo.IsWordBreakSymbol(keyRune) {
		e.ignorePreedit = false
		var processedStr = e.preeditor.GetProcessedString(bamboo.VietnameseMode, true)
		if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(processedStr) {
			processedStr = e.macroTable.GetText(processedStr)
			e.commitText(processedStr + string(keyRune))
			e.resetPreedit()
			return true, nil
		}
		e.commitText(e.getComposedString() + string(keyRune))
		e.resetPreedit()
		return true, nil
	}
	e.commitPreedit()
	return false, nil
}

func (e *IBusBambooEngine) updatePreedit() {
	var processedStr = e.getPreeditString()
	var preeditLen = uint32(len([]rune(processedStr)))
	if preeditLen == 0 {
		e.HidePreeditText()
		return
	}
	var ibusText = ibus.NewText(e.encodeText(processedStr))

	if e.config.IBflags&IBpreeditInvisibility != 0 {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_NONE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	} else {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	}
	e.UpdatePreeditTextWithMode(ibusText, preeditLen, true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)

	mouseCaptureUnlock()
}

func (e *IBusBambooEngine) shouldFallbackToEnglish() bool {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return false
	}
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode, true)
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
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode, true)
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
	return e.preeditor.GetSpellingMatchResult(bamboo.ToneLess, deepSearch)
}

func (e *IBusBambooEngine) getComposedString() string {
	if e.mustFallbackToEnglish() {
		return e.getProcessedString(bamboo.EnglishMode, false)
	}
	return e.getProcessedString(bamboo.VietnameseMode, false)
}

func (e *IBusBambooEngine) encodeText(text string) string {
	return bamboo.Encode(e.config.OutputCharset, text)
}

func (e *IBusBambooEngine) getProcessedString(mode bamboo.Mode, letterOnly bool) string {
	return e.preeditor.GetProcessedString(mode, letterOnly)
}

func (e *IBusBambooEngine) getPreeditString() string {
	if e.shouldFallbackToEnglish() {
		return e.getProcessedString(bamboo.EnglishMode, false)
	}
	return e.getProcessedString(bamboo.VietnameseMode, false)
}

func (e *IBusBambooEngine) getMode() bamboo.Mode {
	if e.shouldFallbackToEnglish() {
		return bamboo.EnglishMode
	}
	return bamboo.VietnameseMode
}

func (e *IBusBambooEngine) resetPreedit() {
	e.HidePreeditText()
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) commitPreedit() {
	e.HidePreeditText()
	e.commitText(e.getComposedString())
	e.resetPreedit()
}

func (e *IBusBambooEngine) commitText(str string) {
	if len(str) == 0 {
		return
	}
	log.Println("Commit Text", str)
	e.CommitText(ibus.NewText(e.encodeText(str)))
}

func (e *IBusBambooEngine) getVnSeq() string {
	return e.preeditor.GetProcessedString(bamboo.VietnameseMode, false)
}

func (e *IBusBambooEngine) hasMacroKey(key string) bool {
	return e.macroTable.GetText(key) != ""
}
