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
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
	"strconv"
)

func (e *IBusBambooEngine) openEmojiList() {
	e.emoji.ProcessKey(':')
	e.UpdatePreeditText(ibus.NewText(":"), 1, true)
	lt := ibus.NewLookupTable()
	lt.Orientation = IBUS_ORIENTATION_HORIZONTAL
	for _, codePoint := range e.emoji.Query() {
		lt.AppendCandidate(codePoint)
	}
	e.emojiLookupTable = lt
	e.emojiUpdateLookupTable()
}

func (e *IBusBambooEngine) emojiProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	var raw = e.emoji.GetRawString()
	var rawTextLen = len([]rune(raw))
	var keyRune = rune(keyVal)
	var cps = e.emoji.Query()
	var reset = func() {
		e.emoji.Reset()
		e.HideAuxiliaryText()
		e.HidePreeditText()
		e.HideLookupTable()
		e.isEmojiTableOpened = false
	}
	if keyRune == ':' {
		reset()
		return false, nil
	}
	if keyVal == IBUS_Return || keyVal == IBUS_KP_Enter {
		if rawTextLen > 0 {
			reset()
			if len(cps) > 0 {
				e.CommitText(ibus.NewText(cps[0]))
			} else {
				e.CommitText(ibus.NewText(raw))
			}
			return true, nil
		}
		return false, nil
	}
	if keyVal == IBUS_Escape {
		if rawTextLen > 0 {
			reset()
			e.CommitText(ibus.NewText(raw))
			return true, nil
		}
		return false, nil
	}
	if keyVal == IBUS_BackSpace {
		if rawTextLen > 0 {
			e.emoji.RemoveLastKey()
		} else {
			return false, nil
		}
	} else if (keyRune >= 'a' && keyRune <= 'z') || (keyRune >= 'A' && keyRune <= 'Z') {
		if raw == ":" {
			e.emoji.Reset()
		}
		e.emoji.ProcessKey(keyRune)
	} else if keyRune >= '1' && keyRune <= '9' {
		if keyNumber, err := strconv.Atoi(string(keyRune)); err == nil {
			e.CommitText(ibus.NewText(cps[keyNumber-1]))
			reset()
			return true, nil
		}
	} else if keyRune > ' ' && keyRune <= '~' {
		e.emoji.ProcessKey(keyRune)
	} else if keyRune < 128 && rawTextLen > 0 {
		reset()
		if len(cps) > 0 {
			e.CommitText(ibus.NewText(cps[0]))
		} else {
			e.CommitText(ibus.NewText(raw))
		}
		return false, nil
	}
	raw = e.emoji.GetRawString()
	rawTextLen = len([]rune(raw))
	cps = e.emoji.Query()
	e.UpdatePreeditTextWithMode(ibus.NewText(raw), uint32(rawTextLen), true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	lt := ibus.NewLookupTable()
	lt.Orientation = IBUS_ORIENTATION_HORIZONTAL
	for _, codePoint := range cps {
		lt.AppendCandidate(codePoint)
	}
	e.emojiLookupTable = lt
	e.emojiUpdateLookupTable()
	mouseCaptureUnlock()
	return true, nil
}

func (e *IBusBambooEngine) emojiUpdateLookupTable() {
	var visible = len(e.emojiLookupTable.Candidates) > 0
	e.UpdateLookupTable(e.emojiLookupTable, visible)
}
