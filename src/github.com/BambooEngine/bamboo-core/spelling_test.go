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

func TestGenerateDumpSoundsFromWord(t *testing.T) {
	s := ParseDumpSoundsFromWord("chuyển")
	if len(s) != 6 {
		t.Errorf("Test length of chuyển, expected [6], got [%v]", len(s))
	}
}

func TestGenerateSoundsFromWord(t *testing.T) {
	s := ParseSoundsFromWord("chuyển")
	if len(s) != 6 {
		t.Errorf("Test length of chuyển, expected [6], got [%v]", len(s))
	}
	s2 := ParseSoundsFromWord("quyển")
	if len(s2) != 5 {
		t.Errorf("Test length of chuyển, expected [5], got [%v]", len(s2))
	}
	s3 := ParseSoundsFromWord("giá")
	if len(s3) != 3 {
		t.Errorf("Test length of chuyển, expected [3], got [%v]", len(s3))
	}
	s4 := ParseSoundsFromWord("gic1")
	if len(s4) != 4 {
		t.Errorf("Test length of chuyển, expected [4], got [%v]", len(s4))
	}
}
