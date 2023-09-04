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
	"time"

	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

func (e *IBusBambooEngine) preeditProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var rawKeyLen = e.getRawKeyLen()
	var keyRune = rune(keyVal)
	var oldText = e.getPreeditString()
	defer e.updateLastKeyWithShift(keyVal, state)

	if !e.shouldRestoreKeyStrokes {
		if !e.preeditor.CanProcessKey(keyRune) && rawKeyLen == 0 && e.config.IBflags&IBmacroEnabled == 0 {
			// don't process special characters if rawKeyLen == 0,
			// workaround for Chrome's address bar and Google SpreadSheets
			return false, nil
		}
	}

	if keyVal == IBusBackSpace {
		if e.runeCount() == 1 {
			e.commitPreeditAndReset("")
			return true, nil
		}
		if rawKeyLen > 0 {
			e.preeditor.RemoveLastChar(true)
			e.updatePreedit(e.getPreeditString())
			return true, nil
		} else {
			return false, nil
		}
	}
	if keyVal == IBusTab {
		if ok, macText := e.getMacroText(); ok {
			e.commitPreeditAndReset(macText)
		} else {
			e.commitPreeditAndReset(e.getComposedString(oldText))
			return false, nil
		}
		return true, nil
	}

	newText, isWordBreakRune := e.getCommitText(keyVal, keyCode, state)
	isValidKey := isValidState(state) && e.isValidKeyVal(keyVal)
	if isWordBreakRune {
		e.commitPreeditAndReset(newText)
		return isValidKey, nil
	}
	e.updatePreedit(newText)
	return isValidKey, nil
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
	defer func() {
		if e.config.IBflags&IBmouseCapturing != 0 {
			mouseCaptureUnlock()
		}
	}()
	var encodedStr = e.encodeText(processedStr)
	var preeditLen = uint32(len([]rune(encodedStr)))
	if preeditLen == 0 {
		e.HidePreeditText()
		e.HideAuxiliaryText()
		e.CommitText(ibus.NewText(""))
		return
	}
	var ibusText = ibus.NewText(encodedStr)
	if inStringList(enabledAuxiliaryTextList, e.getWmClass()) {
		e.UpdateAuxiliaryText(ibusText, true)
		return
	}

	if e.config.IBflags&IBnoUnderline != 0 {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_NONE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	} else {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	}
	e.UpdatePreeditTextWithMode(ibusText, preeditLen, true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
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
	if ok, _ := e.getMacroText(); ok {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && !bamboo.HasAnyVietnameseVower(vnSeq) &&
		(vnRunes[len(vnRunes)-1] == 'd' || strings.ContainsRune(vnSeq, 'đ')) {
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
	if e.config.IBflags&IBmacroEnabled != 0 {
		return e.getProcessedString(bamboo.PunctuationMode)
	}
	if e.shouldFallbackToEnglish(true) {
		return e.getProcessedString(bamboo.EnglishMode)
	}
	return e.getProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) resetPreedit() {
	e.HidePreeditText()
	e.HideAuxiliaryText()
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) commitPreeditAndReset(s string) {
	e.commitText(s)
	e.HidePreeditText()
	e.HideAuxiliaryText()
	e.HideLookupTable()
	e.preeditor.Reset()
}

func (e *IBusBambooEngine) commitText(str string) {
	if str == "" {
		return
	}
	log.Printf("Commit Text [%s]\n", str)
	var now = time.Now()
	e.lastCommitText = now.UnixNano()
	e.CommitText(ibus.NewText(e.encodeText(str)))
}

func (e *IBusBambooEngine) getVnSeq() string {
	return e.preeditor.GetProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) hasMacroKey(key string) bool {
	return e.macroTable.GetText(key) != ""
}
