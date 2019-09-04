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
	"testing"
)

func TestEmojiFindResult(t *testing.T) {
	loadEmojiOne("../../" + DictEmojiOne)
	var be = NewEmojiEngine()
	if be.MatchString(":'") != true {
		t.Errorf("Finding result for emoji :', expected true, got %v", be.MatchString(":'"))
	}
	if be.MatchString(":')") != true {
		t.Errorf("Finding result for emoji :'), expected true, got %v", be.MatchString(":')"))
	}
	if be.MatchString("gri") != true {
		t.Errorf("Finding result for emoji gri, expected true, got %v", be.MatchString("gri"))
	}
	if be.MatchString("grin") != true {
		t.Errorf("Finding result for emoji grin, expected true, got %v", be.MatchString("grinning"))
	}
}

func TestFilterEmoji(t *testing.T) {
	loadEmojiOne("../../" + DictEmojiOne)
	var be = NewEmojiEngine()
	var grinnings = be.Filter(":')")
	if !inStringList(grinnings, "😂") {
		t.Errorf("Filtering emojo :'), expected %v, got %v", true, inStringList(grinnings, "😂"))
	}
	var grinnings2 = be.Filter(":")
	if !inStringList(grinnings2, "😂") {
		t.Errorf("Filtering emojo :, expected %v, got %v", true, inStringList(grinnings2, "😂"))
	}
	var grinnings3 = be.Filter("grin")
	if !inStringList(grinnings3, "😀") {
		t.Errorf("Filtering emojo `grin`, expected %v, got %v", true, inStringList(grinnings3, "😀"))
	}
}
