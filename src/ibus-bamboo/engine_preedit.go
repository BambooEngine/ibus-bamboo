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
	defer e.updateLastKeyWithShift(keyVal, state)

	if !e.isValidState(state) {
		e.commitPreedit(e.getPreeditString())
		return false, nil
	}

	if keyVal == IBUS_BackSpace {
		if rawKeyLen > 0 {
			e.preeditor.RemoveLastChar(true)
			e.updatePreedit(e.getPreeditString())
			return true, nil
		} else {
			return false, nil
		}
	}
	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		var oMode = bamboo.VietnameseMode
		if e.shouldFallbackToEnglish() {
			oMode = bamboo.EnglishMode
		}
		e.preeditor.ProcessKey(keyRune, e.getInputMode())
		var newSeq = e.preeditor.GetProcessedString(oMode | bamboo.WithEffectKeys)
		if len(newSeq) > 0 && rune(newSeq[len(newSeq)-1]) == keyRune && bamboo.IsWordBreakSymbol(keyRune) {
			e.commitPreedit(newSeq)
		} else {
			e.updatePreedit(e.getPreeditString())
		}
		return true, nil
	} else if bamboo.IsWordBreakSymbol(keyRune) || ('0' <= keyVal && keyVal <= '9') {
		if keyVal == IBUS_Space && state&IBUS_SHIFT_MASK != 0 &&
			e.config.IBflags&IBrestoreKeyStrokesEnabled != 0 && !e.lastKeyWithShift {
			// restore key strokes
			var vnSeq = e.preeditor.GetProcessedString(bamboo.VietnameseMode)
			if bamboo.HasAnyVietnameseRune(vnSeq) {
				e.preeditor.RestoreLastWord()
				e.updatePreedit(e.getPreeditString())
			} else {
				e.commitPreedit(vnSeq + string(keyRune))
			}
			return true, nil
		}
		var processedStr = e.preeditor.GetProcessedString(bamboo.VietnameseMode)
		if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(processedStr) {
			processedStr = e.expandMacro(processedStr)
			e.commitPreedit(processedStr + string(keyRune))
			return true, nil
		}
		e.commitPreedit(e.getComposedString() + string(keyRune))
		return true, nil
	}
	e.commitPreedit(e.getPreeditString())
	return false, nil
}

func (e *IBusBambooEngine) expandMacro(str string) string {
	var macroText = e.macroTable.GetText(str)
	if e.config.IBflags&IBautoCapitalizeMacro != 0 {
		switch determineMacroCase(str) {
		case VnCaseAllSmall:
			return strings.ToLower(macroText)
		case VnCaseAllCapital:
			return strings.ToUpper(macroText)
		}
	}
	return macroText
}

func (e *IBusBambooEngine) updatePreedit(processedStr string) {
	var encodedStr = e.encodeText(processedStr)
	var preeditLen = uint32(len([]rune(encodedStr)))
	if preeditLen == 0 {
		e.HidePreeditText()
		e.CommitText(ibus.NewText(""))
		return
	}
	var ibusText = ibus.NewText(encodedStr)

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
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode | bamboo.LowerCase)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return false
	}
	if e.config.IBflags&IBmarcoEnabled != 0 && e.macroTable.HasKey(vnSeq) {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && (vnRunes[len(vnRunes)-1] == 'd' || strings.ContainsRune(vnSeq, 'đ')) {
		if !bamboo.HasAnyVowel(vnRunes) {
			return false
		}
	}
	if !bamboo.HasAnyVietnameseRune(vnSeq) {
		return false
	}
	if e.preeditor.GetSpellingMatchResult(0) != bamboo.FindResultNotMatch {
		return false
	}
	return true
}

func (e *IBusBambooEngine) hasVietnameseChar() bool {
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode)
	return bamboo.HasAnyVietnameseRune(vnSeq)
}

func (e *IBusBambooEngine) mustFallbackToEnglish() bool {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return false
	}
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode | bamboo.LowerCase)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && strings.ContainsRune(vnSeq, 'đ') {
		if !bamboo.HasAnyVowel(vnRunes) {
			return false
		}
	}
	if !bamboo.HasAnyVietnameseRune(vnSeq) {
		return false
	}
	if e.preeditor.GetSpellingMatchResult(bamboo.WithDictionary) == bamboo.FindResultMatchFull {
		return false
	}
	if e.config.IBflags&IBspellCheckingWithDicts == 0 {
		return e.preeditor.GetSpellingMatchResult(0) != bamboo.FindResultMatchFull
	}
	return true
}

func (e *IBusBambooEngine) getComposedString() string {
	if e.mustFallbackToEnglish() {
		return e.getProcessedString(bamboo.EnglishMode | bamboo.WithEffectKeys)
	}
	return e.getProcessedString(bamboo.VietnameseMode | bamboo.WithEffectKeys)
}

func (e *IBusBambooEngine) encodeText(text string) string {
	return bamboo.Encode(e.config.OutputCharset, text)
}

func (e *IBusBambooEngine) getProcessedString(mode bamboo.Mode) string {
	return e.preeditor.GetProcessedString(mode)
}

func (e *IBusBambooEngine) getPreeditString() string {
	if e.shouldFallbackToEnglish() {
		return e.getProcessedString(bamboo.EnglishMode | bamboo.WithEffectKeys)
	}
	return e.getProcessedString(bamboo.VietnameseMode | bamboo.WithEffectKeys)
}

func (e *IBusBambooEngine) getInputMode() bamboo.Mode {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return bamboo.VietnameseMode
	}
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode | bamboo.LowerCase)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return bamboo.VietnameseMode
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && (vnRunes[len(vnRunes)-1] == 'd' || strings.ContainsRune(vnSeq, 'đ')) {
		if !bamboo.HasAnyVowel(vnRunes) {
			return bamboo.VietnameseMode
		}
	}
	if e.preeditor.GetSpellingMatchResult(0) != bamboo.FindResultNotMatch {
		return bamboo.VietnameseMode
	}
	return bamboo.EnglishMode
}

func (e *IBusBambooEngine) resetPreedit() {
	e.HidePreeditText()
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) commitPreedit(s string) {
	e.HidePreeditText()
	e.commitText(s)
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) commitText(str string) {
	if str == "" {
		return
	}
	log.Println("Commit Text", str)
	e.CommitText(ibus.NewText(e.encodeText(str)))
}

func (e *IBusBambooEngine) getVnSeq() string {
	return e.preeditor.GetProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) hasMacroKey(key string) bool {
	return e.macroTable.GetText(key) != ""
}
