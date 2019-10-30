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
	"strconv"
)

const EMOJI_MAX_PAGE_SIZE = 9

func (e *IBusBambooEngine) openEmojiList() {
	e.emoji.ProcessKey(':')
	e.UpdatePreeditText(ibus.NewText(":"), 1, true)
	e.UpdateAuxiliaryText(ibus.NewText(":"), true)
	lt := ibus.NewLookupTable()
	lt.Orientation = IBUS_ORIENTATION_HORIZONTAL
	for _, codePoint := range e.emoji.Query() {
		lt.AppendCandidate(codePoint)
	}
	lt.PageSize = uint32(EMOJI_MAX_PAGE_SIZE)
	e.emojiLookupTable = lt
	e.updateEmojiLookupTable()
}

func (e *IBusBambooEngine) emojiProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var raw = e.emoji.GetRawString()
	var rawTextLen = len([]rune(raw))
	var keyRune = rune(keyVal)
	var reset = e.closeEmojiCandidates
	if keyVal == IBUS_Colon {
		reset()
		return false, nil
	}
	if keyVal == IBUS_Return {
		if rawTextLen > 0 {
			if len(e.emojiLookupTable.Candidates) > 0 {
				e.commitEmojiCandidate()
			} else {
				e.CommitText(ibus.NewText(raw))
			}
			reset()
			return true, nil
		}
		return false, nil
	}
	if keyVal == IBUS_Escape {
		if rawTextLen > 0 {
			e.CommitText(ibus.NewText(raw))
			reset()
			return true, nil
		}
		return false, nil
	}
	if keyVal == IBUS_Left || keyVal == IBUS_Up {
		e.CursorUp()
		return true, nil
	} else if keyVal == IBUS_Right || keyVal == IBUS_Down {
		e.CursorDown()
		return true, nil
	} else if keyVal == IBUS_Page_Up {
		e.PageUp()
		return true, nil
	} else if keyVal == IBUS_Page_Down {
		e.PageDown()
		return true, nil
	}
	if keyVal == IBUS_BackSpace {
		if rawTextLen > 0 {
			e.emoji.RemoveLastKey()
		} else {
			reset()
			return false, nil
		}
	} else if (keyRune >= 'a' && keyRune <= 'z') || (keyRune >= 'A' && keyRune <= 'Z') {
		var testStr = string(append(e.emoji.keys, keyRune))
		if raw == ":" && e.emoji.MatchString(testStr) == false {
			e.emoji.Reset()
		}
		e.emoji.ProcessKey(keyRune)
	} else if keyRune >= '1' && keyRune <= '9' {
		if pos, err := strconv.Atoi(string(keyRune)); err == nil {
			if e.setCursorPosInEmojiTable(uint32(pos - 1)) {
				e.commitEmojiCandidate()
				reset()
				return true, nil
			} else {
				reset()
			}
		}
		return false, nil
	} else if (keyRune >= ' ' && keyRune <= '~') || bamboo.IsWordBreakSymbol(keyRune) {
		var testStr = string(append(e.emoji.keys, keyRune))
		if raw == ":" && e.emoji.MatchString(testStr) == false {
			e.emoji.Reset()
		}
		e.emoji.ProcessKey(keyRune)
		if e.emoji.MatchString(string(e.emoji.keys)) == false {
			e.CommitText(ibus.NewText(e.emoji.GetRawString()))
			reset()
			return true, nil
		}
	} else if rawTextLen > 0 {
		reset()
		e.CommitText(ibus.NewText(raw))
		return false, nil
	}
	raw = e.emoji.GetRawString()
	rawTextLen = len([]rune(raw))
	cps := e.emoji.Query()
	if cps != nil {
		codePoint0 := cps[0]
		e.UpdatePreeditTextWithMode(ibus.NewText(codePoint0), uint32(len(codePoint0)), true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	} else {
		e.UpdatePreeditTextWithMode(ibus.NewText(raw), uint32(rawTextLen), true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	}
	e.UpdateAuxiliaryText(ibus.NewText(raw), true)
	lt := ibus.NewLookupTable()
	lt.Orientation = IBUS_ORIENTATION_HORIZONTAL
	for _, codePoint := range cps {
		lt.AppendCandidate(codePoint)
	}
	lt.PageSize = uint32(EMOJI_MAX_PAGE_SIZE)
	e.emojiLookupTable = lt
	e.updateEmojiLookupTable()
	return true, nil
}

func (e *IBusBambooEngine) setCursorPosInEmojiTable(idx uint32) bool {
	pageSize := e.emojiLookupTable.PageSize
	if idx > pageSize {
		return false
	}
	page := e.emojiLookupTable.CursorPos / pageSize
	newPos := page*pageSize + idx
	if int(newPos) > len(e.emojiLookupTable.Candidates) {
		return false
	}
	e.emojiLookupTable.CursorPos = newPos
	return true
}

func (e *IBusBambooEngine) updateEmojiLookupTable() {
	var visible = len(e.emojiLookupTable.Candidates) > 0
	e.UpdateLookupTable(e.emojiLookupTable, visible)
	var cps = e.emoji.Query()
	if pos := e.emojiLookupTable.CursorPos; pos < uint32(len(cps)) {
		var codePoint0 = cps[pos]
		e.UpdatePreeditTextWithMode(ibus.NewText(codePoint0), uint32(len(codePoint0)), true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	}
}

func (e *IBusBambooEngine) commitEmojiCandidate() {
	var cps = e.emoji.Query()
	if pos := e.emojiLookupTable.CursorPos; pos < uint32(len(cps)) {
		e.CommitText(ibus.NewText(cps[pos]))
	}
}

func (e *IBusBambooEngine) refreshEmojiCandidate() {
	var raw = e.emoji.GetRawString()
	var rawTextLen = len([]rune(raw))
	e.UpdatePreeditTextWithMode(ibus.NewText(raw), uint32(rawTextLen), true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	e.UpdateAuxiliaryText(ibus.NewText(raw), true)
	e.updateEmojiLookupTable()
}

func (e *IBusBambooEngine) closeEmojiCandidates() {
	e.emojiLookupTable = nil
	e.emoji.Reset()
	e.UpdateLookupTable(ibus.NewLookupTable(), true) // workaround for issue #18
	e.HidePreeditText()
	e.HideLookupTable()
	e.HideAuxiliaryText()
	e.isEmojiLTOpened = false
}
