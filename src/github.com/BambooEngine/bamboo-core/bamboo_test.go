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

func newStdEngine() IEngine {
	return NewEngine("Telex 2", EstdFlags, nil)
}

func TestProcessString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("aw", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ă" {
		t.Errorf("Process [aw], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "ă")
	}
	ng.Reset()
	ng.ProcessString("uwow", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ươ" {
		t.Errorf("Process [uwow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "ươ")
	}
	ng.Reset()
	ng.ProcessString("chuaarn", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "chuẩn" {
		t.Errorf("Process [chuaarn], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "chuẩn")
	}
	ng.Reset()
	ng.ProcessString("giamaf", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "giầm" {
		t.Errorf("Process [giamaf], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "giầm")
	}
}

func TestProcessDDString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("dd", VietnameseMode)
	if ng.IsSpellingCorrect(NoTone) != true {
		t.Errorf("IsSpellingCorrect [dd], got [%v] expected [true]", ng.IsSpellingCorrect(NoTone))
	}
	ng.Reset()
	ng.ProcessString("ddafi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "đài" {
		t.Errorf("Process [ddafi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "đài")
	}
}

func TestProcessMuoiwqString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("Muoiwq", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "Muoiwq" {
		t.Errorf("Process [Muoiwq], got [%s] expected [Muoiwq]", ng.GetProcessedString(EnglishMode, false))
	}
	ng.Reset()
	ng.ProcessString("mootj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "một" {
		t.Errorf("Process [mootj], got [%s] expected [một]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessThuowString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("Thuow", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "Thuơ" {
		t.Errorf("Process [Thuow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "Thuơ")
	}
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode, false) != "Thu" {
		t.Errorf("Process [Thuow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "Thu")
	}
}

func TestProcessUpperString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("VIEETJ", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "VIỆT" {
		t.Errorf("Process [VIEETJ], got [%s] expected [VIỆT]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode, false) != "VIỆ" {
		t.Errorf("Process remove last char of upper string, got [%s] expected [VIỆ]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.ProcessChar('Q', VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "VIEEJQ" {
		t.Errorf("Process remove last char of upper string, got [%s] expected [VIEEJQ]", ng.GetProcessedString(EnglishMode, false))
	}
	ng.Reset()
	ng.ProcessString("IB", EnglishMode)
	if ng.GetProcessedString(EnglishMode, false) != "IB" {
		t.Errorf("Process remove last char of upper string, got [%s] expected [IB]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestSpellingCheck(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("noww", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "now" {
		t.Errorf("Process [noww], got [%s] expected [now]", ng.GetProcessedString(EnglishMode, false))
	}
	ng.Reset()
	ng.ProcessString("sawss", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "saws" {
		t.Errorf("Process [sawss], got [%s] expected [saws]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestProcessDD(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("dd", VietnameseMode)
	if ng.IsSpellingCorrect(NoTone) != true {
		t.Errorf("Check spelling for [dd], got [%v] expected [true]", ng.IsSpellingCorrect(NoTone))
	}
	if ng.GetProcessedString(VietnameseMode, false) != "đ" {
		t.Errorf("Process [dd], got [%s] expected [đ]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestTelex3(t *testing.T) {
	ng := NewEngine("Telex 3", EstdFlags, nil)
	ng.ProcessString("[", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ươ" {
		t.Errorf("Process Telex 3 [[], got [%v] expected [ươ]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	ng.ProcessString("{", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ƯƠ" {
		t.Errorf("Process Telex 3 [{], got [%s] expected [ƯƠ]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestProcessNguwowfiString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("wowfi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ười" {
		t.Errorf("Process [wowfi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "ười")
	}
}

func TestRemoveLastChar(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("hanhj", VietnameseMode)
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode, false) != "hạn" {
		t.Errorf("Process [hanhj], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "hạn")
	}
	ng.Reset()
}

func TestProcessCatrString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("catr", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "catr" {
		t.Errorf("Process [nguwowfi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "catr")
	}
}

func TestProcessToowiString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("toowi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tơi" {
		t.Errorf("Process [toowi], got [%s] expected [tơi]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessAlooString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("aloo", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "alô" {
		t.Errorf("Process [aloo], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "alô")
	}
}

func TestSpellingCheckForGiw(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("giw", VietnameseMode)
	if ng.IsSpellingCorrect(NoTone|NoMark) != true {
		t.Errorf("TestSpellingCheckForGiw, got [%v] expected [%v]", ng.IsSpellingCorrect(NoTone|NoMark), true)
	}
}

func TestDoubleBrackets(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("[[", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "[" {
		t.Errorf("TestDoubleBrackets, got [%v] expected [%v]", ng.GetProcessedString(EnglishMode, false), "[")
	}
}
func TestDoubleBracketso(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("tooss", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tôs" {
		t.Errorf("TestDoubleBrackets tooss, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	ng.ProcessString("tosos", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tôs" {
		t.Errorf("TestDoubleBrackets tosos, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestDoubleW(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("ww", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "w" {
		t.Errorf("TestDoubleW, got [%v] expected [w]", ng.GetProcessedString(EnglishMode, false))
	}
	if ng.GetProcessedString(VietnameseMode, false) != "w" {
		t.Errorf("TestDoubleW, got [%v] expected [w]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestDoubleW2(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("wiw", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "uiw" {
		t.Errorf("TestDoubleW, got [%v] expected [uiw]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestProcessDuwoi(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("duwoi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "dươi" {
		t.Errorf("TestProcessDuwoi, got [%v] expected [dươi]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessRefresh(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("reff", VietnameseMode)
	ng.ProcessString("resh", EnglishMode)
	if ng.GetProcessedString(EnglishMode, false) != "refresh" {
		t.Errorf("TestProcessDuwoi, got [%v] expected [refresh]", ng.GetProcessedString(VietnameseMode, false))
	}
}
func TestProcessRefresh2(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("reff", VietnameseMode)
	ng.RemoveLastChar()
	ng.ProcessChar('f', VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "rè" {
		t.Errorf("TestProcessDuwoi, got [%v] expected [rè]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessDDSeq(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("oddp", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ođp" {
		t.Errorf("TestProcessDDSeq, got [%v] expected [ođp]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessGisa(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("gisa", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "giá" {
		t.Errorf("TestProcessDDSeq, got [%v] expected [giá]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessKimso(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("kimso", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "kímo" {
		t.Errorf("TestProcessKimso, got [%v] expected [kímo]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessTo(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("to", VietnameseMode)
	if ng.IsSpellingCorrect(VietnameseMode|NoTone) != true {
		t.Errorf("TestProcessTo, got [%v] expected [true]", ng.IsSpellingCorrect(VietnameseMode))
	}
}

func TestProcessToorr(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("toorr", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tôr" {
		t.Errorf("TestProcessToorr, got [%v] expected [tôr]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//tnó
func TestProcessTnoss(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("tnoss", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tnos" {
		t.Errorf("TestProcessToorr, got [%v] expected [tnos]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//ềng
func TestProcessEenghf(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags, map[string]bool{"ềngh": true})
	ng.ProcessString("eenghf", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ềngh" {
		t.Errorf("TestProcessToorr, got [%v] expected [ềnhg]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//HIEEUR
func TestProcessHIEEUR(t *testing.T) {
	ng := NewEngine("Telex 2", EstdFlags, nil)
	ng.ProcessString("tooi oo HIEEUR", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tôi ô HIỂU" {
		t.Errorf("TestProcessToorr, got [%v] expected [tôi ô HIỂU]", ng.GetProcessedString(VietnameseMode, false))
	}
	if ng.GetProcessedString(VietnameseMode, true) != "HIỂU" {
		t.Errorf("TestProcessToorr, got [%v] expected [HIỂU]", ng.GetProcessedString(VietnameseMode, true))
	}
}
