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

	"ibus-bamboo/emoji"
)

func TestEmojiFindResult(t *testing.T) {
	emoji.Trie, _ = emoji.LoadEmojiOne(DictEmojiOne)
	var be = emoji.NewEmojiEngine()
	if emoji.MatchString(be, ":'") != true {
		t.Errorf("Finding result for emoji :', expected true, got %v", emoji.MatchString(be, ":'"))
	}
	if emoji.MatchString(be, ":')") != true {
		t.Errorf("Finding result for emoji :'), expected true, got %v", emoji.MatchString(be, ":')"))
	}
	if emoji.MatchString(be, "gri") != true {
		t.Errorf("Finding result for emoji gri, expected true, got %v", emoji.MatchString(be, "gri"))
	}
	if emoji.MatchString(be, "grin") != true {
		t.Errorf("Finding result for emoji grin, expected true, got %v", emoji.MatchString(be, "grinning"))
	}
}

func TestFilterEmoji(t *testing.T) {
	emoji.Trie, _ = emoji.LoadEmojiOne(DictEmojiOne)
	var be = emoji.NewEmojiEngine()
	var grinnings = emoji.Filter(be, ":')")
	if !inStringList(grinnings, "😂") {
		t.Errorf("Filtering emojo :'), expected %v, got %v", true, inStringList(grinnings, "😂"))
	}
	var grinnings2 = emoji.Filter(be, ":")
	if !inStringList(grinnings2, "😂") {
		t.Errorf("Filtering emojo :, expected %v, got %v", true, inStringList(grinnings2, "😂"))
	}
	var grinnings3 = emoji.Filter(be, "grin")
	if !inStringList(grinnings3, "😀") {
		t.Errorf("Filtering emojo `grin`, expected %v got %v", true, inStringList(grinnings3, "😀"))
	}
}
