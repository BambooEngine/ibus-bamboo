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
	"log"
	"strings"

	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

func (e *IBusBambooEngine) preeditProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var rawKeyLen = e.getRawKeyLen()
	var keyRune = rune(keyVal)
	var oldText = e.getPreeditString()
	defer e.updateLastKeyWithShift(keyVal, state)

	// workaround for chrome's address bar and Google SpreadSheets
	if !e.isValidState(state) || !e.canProcessKey(keyVal) ||
		(rawKeyLen == 0 && !e.preeditor.CanProcessKey(keyRune)) {
		if rawKeyLen > 0 {
			e.HidePreeditText()
			e.commitText(e.getPreeditString())
			e.preeditor.Reset()
		}
		return false, nil
	}

	if keyVal == IBusBackSpace {
		if rawKeyLen > 0 {
			e.preeditor.RemoveLastChar(true)
			e.updatePreedit(e.getPreeditString())
			return true, nil
		} else {
			return false, nil
		}
	}
	if keyVal == IBusTab {
		var text = e.preeditor.GetProcessedString(bamboo.VietnameseMode)
		if e.config.IBflags&IBmacroEnabled != 0 && e.macroTable.HasKey(text) {
			// macro processing
			macText := e.expandMacro(text)
			e.commitPreedit(macText)
		} else {
			e.commitPreedit(e.getComposedString(oldText))
			return false, nil
		}
		return true, nil
	}

	if e.preeditor.CanProcessKey(keyRune) {
		if state&IBusLockMask != 0 {
			keyRune = e.toUpper(keyRune)
		}
		e.preeditor.ProcessKey(keyRune, e.getBambooInputMode())
		if inKeyList(e.preeditor.GetInputMethod().AppendingKeys, keyRune) {
			if fullSeq := e.preeditor.GetProcessedString(bamboo.VietnameseMode); len(fullSeq) > 0 && rune(fullSeq[len(fullSeq)-1]) == keyRune {
				e.commitPreedit(fullSeq)
			} else if newText := e.getPreeditString(); newText != "" && keyRune == rune(newText[len(newText)-1]) {
				e.commitPreedit(oldText + string(keyRune))
			} else {
				e.updatePreedit(e.getPreeditString())
			}
		} else {
			e.updatePreedit(e.getPreeditString())
		}
		return true, nil
	} else if bamboo.IsWordBreakSymbol(keyRune) {
		if keyVal == IBusSpace && state&IBusShiftMask != 0 &&
			e.config.IBflags&IBrestoreKeyStrokesEnabled != 0 && !e.lastKeyWithShift {
			// restore key strokes
			var vnSeq = e.preeditor.GetProcessedString(bamboo.VietnameseMode)
			if bamboo.HasAnyVietnameseRune(vnSeq) {
				e.commitPreedit(e.preeditor.GetProcessedString(bamboo.EnglishMode))
			} else {
				e.commitPreedit(vnSeq + string(keyRune))
			}
			return true, nil
		}
		var processedStr = e.preeditor.GetProcessedString(bamboo.VietnameseMode)
		if e.config.IBflags&IBmacroEnabled != 0 && e.macroTable.HasKey(processedStr) {
			processedStr = e.expandMacro(processedStr)
			e.commitPreedit(processedStr + string(keyRune))
			return true, nil
		}
		e.commitPreedit(e.getComposedString(oldText) + string(keyRune))
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

	if e.config.IBflags&IBmouseCapturing != 0 {
		mouseCaptureUnlock()
	}
}

func (e *IBusBambooEngine) getWhiteList() [][]string {
	return [][]string{
		e.config.PreeditWhiteList,
		e.config.SurroundingTextWhiteList,
		e.config.ForwardKeyWhiteList,
		e.config.SLForwardKeyWhiteList,
		e.config.X11ClipboardWhiteList,
		e.config.DirectForwardKeyWhiteList,
		e.config.ExceptedList,
	}
}

func (e *IBusBambooEngine) getBambooInputMode() bamboo.Mode {
	if e.shouldFallbackToEnglish(false) {
		return bamboo.EnglishMode
	}
	return bamboo.VietnameseMode
}

func (e *IBusBambooEngine) shouldFallbackToEnglish(checkVnRune bool) bool {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return false
	}
	var vnSeq = e.getProcessedString(bamboo.VietnameseMode | bamboo.LowerCase)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return false
	}
	if e.config.IBflags&IBmacroEnabled != 0 && e.macroTable.HasKey(vnSeq) {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && (vnRunes[len(vnRunes)-1] == 'd' || strings.ContainsRune(vnSeq, 'đ')) {
		return false
	}
	if checkVnRune && !bamboo.HasAnyVietnameseRune(vnSeq) {
		return false
	}
	return !e.preeditor.IsValid(false)
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
		return false
	}
	if e.config.IBflags&IBspellCheckWithDicts != 0 {
		return !dictionary[vnSeq]
	}
	return !e.preeditor.IsValid(true)
}

func (e *IBusBambooEngine) getComposedString(oldText string) string {
	if bamboo.HasAnyVietnameseRune(oldText) && e.mustFallbackToEnglish() {
		return e.getProcessedString(bamboo.EnglishMode)
	}
	return oldText
}

func (e *IBusBambooEngine) encodeText(text string) string {
	return bamboo.Encode(e.config.OutputCharset, text)
}

func (e *IBusBambooEngine) getProcessedString(mode bamboo.Mode) string {
	return e.preeditor.GetProcessedString(mode)
}

func (e *IBusBambooEngine) getPreeditString() string {
	if e.shouldFallbackToEnglish(true) {
		return e.getProcessedString(bamboo.EnglishMode)
	}
	return e.getProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) resetPreedit() {
	e.HidePreeditText()
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) commitPreedit(s string) {
	e.commitText(s)
	e.HidePreeditText()
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) commitText(str string) {
	if str == "" {
		return
	}
	log.Printf("Commit Text [%s]\n", str)
	e.CommitText(ibus.NewText(e.encodeText(str)))
}

func (e *IBusBambooEngine) getVnSeq() string {
	return e.preeditor.GetProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) hasMacroKey(key string) bool {
	return e.macroTable.GetText(key) != ""
}
