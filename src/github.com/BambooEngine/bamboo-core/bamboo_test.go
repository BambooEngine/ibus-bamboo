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

func TestProcessString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("aw")
	if ng.GetProcessedString(VietnameseMode) != "ă" {
		t.Errorf("Process [aw], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "ă")
	}
	ng.Reset()
	ng.ProcessString("uwow")
	if ng.GetProcessedString(VietnameseMode) != "ươ" {
		t.Errorf("Process [uwow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "ươ")
	}
	ng.Reset()
	ng.ProcessString("chuaarn")
	if ng.GetProcessedString(VietnameseMode) != "chuẩn" {
		t.Errorf("Process [chuaarn], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "chuẩn")
	}
	ng.Reset()
	ng.ProcessString("giamaf")
	if ng.GetProcessedString(VietnameseMode) != "giầm" {
		t.Errorf("Process [giamaf], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "giầm")
	}
}

func TestProcessDDString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("dd")
	if ng.IsSpellingCorrect(NoTone) != true {
		t.Errorf("IsSpellingCorrect [dd], got [%v] expected [true]", ng.IsSpellingCorrect(NoTone))
	}
	ng.Reset()
	ng.ProcessString("ddafi")
	if ng.GetProcessedString(VietnameseMode) != "đài" {
		t.Errorf("Process [ddafi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "đài")
	}
}

func TestProcessMuoiwqString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("Muoiwq")
	if ng.GetProcessedString(EnglishMode) != "Muoiwq" {
		t.Errorf("Process [Muoiwq], got [%s] expected [Muoiwq]", ng.GetProcessedString(EnglishMode))
	}
}

func TestProcessThuowString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("Thuow")
	if ng.GetProcessedString(VietnameseMode) != "Thuơ" {
		t.Errorf("Process [Thuow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "Thuơ")
	}
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode) != "Thu" {
		t.Errorf("Process [Thuow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "Thu")
	}
}

func TestProcessUpperString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("VIEETJ")
	if ng.GetProcessedString(VietnameseMode) != "VIỆT" {
		t.Errorf("Process [VIEETJ], got [%s] expected [VIỆT]", ng.GetProcessedString(VietnameseMode))
	}
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode) != "VIỆ" {
		t.Errorf("Process remove last char of upper string, got [%s] expected [VIỆ]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestSpellingCheck(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("noww")
	if ng.GetProcessedString(EnglishMode) != "now" {
		t.Errorf("Process [noww], got [%s] expected [now]", ng.GetProcessedString(EnglishMode))
	}
	ng.Reset()
	ng.ProcessString("sawss")
	if ng.GetProcessedString(EnglishMode) != "saws" {
		t.Errorf("Process [sawss], got [%s] expected [saws]", ng.GetProcessedString(EnglishMode))
	}
}

func TestProcessDD(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("dd")
	if ng.IsSpellingCorrect(NoTone) != true {
		t.Errorf("Check spelling for [dd], got [%v] expected [true]", ng.IsSpellingCorrect(NoTone))
	}
	if ng.GetProcessedString(VietnameseMode) != "đ" {
		t.Errorf("Process [dd], got [%s] expected [đ]", ng.GetProcessedString(EnglishMode))
	}
}

func TestTelex3(t *testing.T) {
	ng := NewEngine("Telex 3", EstdFlags)
	ng.ProcessString("[")
	if ng.GetProcessedString(VietnameseMode) != "ươ" {
		t.Errorf("Process Telex 3 [[], got [%v] expected [ươ]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("{")
	if ng.GetProcessedString(VietnameseMode) != "ƯƠ" {
		t.Errorf("Process Telex 3 [{], got [%s] expected [ƯƠ]", ng.GetProcessedString(EnglishMode))
	}
}

func TestProcessNguwowfiString(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("wowfi")
	if ng.GetProcessedString(VietnameseMode) != "ười" {
		t.Errorf("Process [wowfi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "ười")
	}
}

func TestRemoveLastChar(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("hanhj")
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode) != "hạn" {
		t.Errorf("Process [hanhj], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "hạn")
	}
	ng.Reset()
}

func TestProcessCatrString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("catr")
	if ng.GetProcessedString(VietnameseMode) != "catr" {
		t.Errorf("Process [nguwowfi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "catr")
	}
}

func TestProcessToowiString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags)
	ng.ProcessString("toowi")
	if ng.GetProcessedString(VietnameseMode) != "tơi" {
		t.Errorf("Process [toowi], got [%s] expected [tơi]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessAlooString(t *testing.T) {
	ng := NewEngine("Telex", EstdFlags&^EspellCheckEnabled)
	ng.ProcessString("aloo")
	if ng.GetProcessedString(VietnameseMode) != "alô" {
		t.Errorf("Process [aloo], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "alô")
	}
}

func TestSpellingCheckForGiw(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("giw")
	if ng.IsSpellingCorrect(NoTone|NoMark) != true {
		t.Errorf("TestSpellingCheckForGiw, got [%v] expected [%v]", ng.IsSpellingCorrect(NoTone|NoMark), true)
	}
}

func TestDoubleBrackets(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("[[")
	if ng.GetProcessedString(EnglishMode) != "[" {
		t.Errorf("TestDoubleBrackets, got [%v] expected [%v]", ng.GetProcessedString(EnglishMode), "[")
	}
}
func TestDoubleBracketso(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("tooss")
	if ng.GetProcessedString(VietnameseMode) != "tôs" {
		t.Errorf("TestDoubleBrackets tooss, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("tosos")
	if ng.GetProcessedString(VietnameseMode) != "tôs" {
		t.Errorf("TestDoubleBrackets tosos, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestDoubleW(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("ww")
	if ng.GetProcessedString(EnglishMode) != "w" {
		t.Errorf("TestDoubleW, got [%v] expected [w]", ng.GetProcessedString(EnglishMode))
	}
	if ng.GetProcessedString(VietnameseMode) != "w" {
		t.Errorf("TestDoubleW, got [%v] expected [w]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestDoubleW2(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("wiw")
	if ng.GetProcessedString(EnglishMode) != "uiw" {
		t.Errorf("TestDoubleW, got [%v] expected [uiw]", ng.GetProcessedString(EnglishMode))
	}
}

func TestProcessDuwoi(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("duwoi")
	if ng.GetProcessedString(VietnameseMode) != "dươi" {
		t.Errorf("TestProcessDuwoi, got [%v] expected [dươi]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessRefresh(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags)
	ng.ProcessString("reffresh")
	if ng.GetProcessedString(EnglishMode) != "refresh" {
		t.Errorf("TestProcessDuwoi, got [%v] expected [refresh]", ng.GetProcessedString(VietnameseMode))
	}
}