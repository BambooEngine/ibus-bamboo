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
	"testing"
)

func TestLoadEmojiData(t *testing.T) {
	data, err := loadEmojiOne("../../" + DictEmojiOne)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 2666 {
		t.Errorf("Emoji data's length is wrong, expected 2666 got %d", len(data))
	}
	t.Log(data["1f602"])
	if data["1f602"].Ascii[0] != ":')" {
		t.Errorf("The Emoji ascii is wrong, expected :'), got %s", data["1f602"].Ascii[0])
	}
	if data["1f602"].Name != "face with tears of joy" {
		t.Errorf("The Emoji name is wrong, expected (face with tears of joy), got (%s)", data["1f602"].Name)
	}
}

func TestEmojiFindResult(t *testing.T) {
	var be = NewBambooEmoji("../../" + DictEmojiOne)
	if be.TestString(":'") != bamboo.FindResultMatchPrefix {
		t.Errorf("Finding result for emoji :', expected %d, got %d", bamboo.FindResultMatchPrefix, be.TestString(":'"))
	}
	if be.TestString(":')") != bamboo.FindResultMatchFull {
		t.Errorf("Finding result for emoji :'), expected %d, got %d", bamboo.FindResultMatchFull, be.TestString(":')"))
	}
	if be.TestString("gri") != bamboo.FindResultMatchPrefix {
		t.Errorf("Finding result for emoji gri, expected %d, got %d", bamboo.FindResultMatchPrefix, be.TestString("gri"))
	}
	if be.TestString("grinning") != bamboo.FindResultMatchFull {
		t.Errorf("Finding result for emoji grinning, expected %d, got %d", bamboo.FindResultMatchFull, be.TestString("grinning"))
	}
}

func TestFilterEmoji(t *testing.T) {
	var be = NewBambooEmoji("../../" + DictEmojiOne)
	var grinnings = be.Filter(":')")
	if !inStringList(grinnings, "😂") {
		t.Errorf("Filtering emojo :'), expected %v, got %v", true, inStringList(grinnings, "😂"))
	}
	var grinnings2 = be.Filter(":")
	if !inStringList(grinnings2, "😂") {
		t.Errorf("Filtering emojo :, expected %v, got %v", true, inStringList(grinnings2, "😂"))
	}
	var grinnings3 = be.Filter("grinning")
	if !inStringList(grinnings3, "😀") {
		t.Errorf("Filtering emojo `grinning`, expected %v, got %v", true, inStringList(grinnings3, "😀"))
	}
}
