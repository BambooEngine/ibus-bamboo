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
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 *
 */

package bamboo

import (
	"testing"
)

func TestIsVowel(t *testing.T) {
	if isVowel('a') == false {
		t.Errorf("a is a vowel, but result is false")
	}
	if isVowel('á') == false {
		t.Errorf("á is a vowel, but result is false")
	}
	if isVowel('b') {
		t.Errorf("b is not a vowel, but result is true")
	}
	tvowels := []rune("aàáảãạăằắẳẵặâầấẩẫậeèéẻẽẹêềếểễệiìíỉĩịoòóỏõọôồốổỗộơờớởỡợuùúủũụưừứửữựyỳýỷỹỵ")
	for _, v := range tvowels {
		if isVowel(v) == false {
			t.Errorf("%c is a vowel, but the result is false", v)
		}
	}
}

func TestGetToneFromChar(t *testing.T) {
	none := FindToneFromChar('e')
	if none != TONE_NONE {
		t.Errorf("Test none tune. Got %d, expected %d", none, TONE_NONE)
	}
	grave := FindToneFromChar('è')
	if grave != TONE_GRAVE {
		t.Errorf("Test grave tune. Got %d, expected %d", grave, TONE_GRAVE)
	}
	acute := FindToneFromChar('é')
	if acute != TONE_ACUTE {
		t.Errorf("Test acute tune. Got %d, expected %d", acute, TONE_ACUTE)
	}
	tilde := FindToneFromChar('ẽ')
	if tilde != TONE_TILDE {
		t.Errorf("Test acute tune. Got %d, expected %d", tilde, TONE_TILDE)
	}
	hook := FindToneFromChar('ẻ')
	if hook != TONE_HOOK {
		t.Errorf("Test hook tune. Got %d, expected %d", hook, TONE_HOOK)
	}
	dot := FindToneFromChar('ạ')
	if dot != TONE_DOT {
		t.Errorf("Test dot tune. Got %d, expected %d", dot, TONE_DOT)
	}
}

func TestAddToneToChar(t *testing.T) {
	c_ạ := AddToneToChar('a', uint8(TONE_DOT))
	if c_ạ != 'ạ' {
		t.Errorf("Test add a dot to char. Got %c, expected %c", c_ạ, 'ạ')
	}
}

func TestAddMarkToChar(t *testing.T) {
	c_ặ := AddMarkToChar('ạ', uint8(MARK_BREVE))
	if c_ặ != 'ặ' {
		t.Errorf("Test add a breve to char. Got %c, expected %c", c_ặ, 'ặ')
	}
}
