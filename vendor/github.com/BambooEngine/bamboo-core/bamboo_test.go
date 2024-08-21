/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENSE>.
 */

package bamboo

import (
	"testing"
)

func newStdEngine() IEngine {
	var im = ParseInputMethod(InputMethodDefinitions, "Telex 2")
	return NewEngine(im, EstdFlags)
}

func TestProcessString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("aw", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "ă" {
		t.Errorf("Process [aw], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "ă")
	}
	ng.Reset()
	ng.ProcessString("uw", VietnameseMode)
	ng.ProcessString("o", VietnameseMode)
	ng.ProcessString("w", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "ươ" {
		t.Errorf("Process [uwow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "ươ")
	}
	ng.Reset()
	ng.ProcessString("chuaarn", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "chuẩn" {
		t.Errorf("Process [chuaarn], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "chuẩn")
	}
	ng.Reset()
	ng.ProcessString("giamaf", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "giầm" {
		t.Errorf("Process [giamaf], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "giầm")
	}
}

func TestProcessDDString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("dd", VietnameseMode)
	if ng.IsValid(false) == false {
		t.Errorf("IsSpellingCorrect [dd], got [%v] expected [true]", ng.IsValid(false) == false)
	}
	ng.Reset()
	ng.ProcessString("ddafi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "đài" {
		t.Errorf("Process [ddafi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "đài")
	}
}

func TestProcessMuoiwqString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("Muoiwq", VietnameseMode)
	if ng.GetProcessedString(EnglishMode) != "Muoiwq" {
		t.Errorf("Process [Muoiwq], got [%s] expected [Muoiwq]", ng.GetProcessedString(EnglishMode))
	}
	ng.Reset()
	ng.ProcessString("mootj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "một" {
		t.Errorf("Process [mootj], got [%s] expected [một]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessThuowString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("Thuow", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "Thuơ" {
		t.Errorf("Process [Thuow], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "Thuơ")
	}
	ng.RemoveLastChar(true)
	if ng.GetProcessedString(VietnameseMode) != "Thu" {
		t.Errorf("Process [Thuow] and remove last char, got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "Thu")
	}
}

func TestBambooEngine_RemoveLastChar(t *testing.T) {
	ng := newStdEngine()
	ng.RemoveLastChar(true)
	ng.ProcessString(" ", EnglishMode)
	ng.RemoveLastChar(true)
	ng.ProcessString("loanj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "loạn" {
		t.Errorf("Process [loanj], got [%s] expected [loạn]", ng.GetProcessedString(VietnameseMode))
	}
	ng.RemoveLastChar(true)
	if ng.GetProcessedString(VietnameseMode) != "lọa" {
		t.Errorf("Process [loanj-1], got [%s] expected [lọa]", ng.GetProcessedString(VietnameseMode))
	}
	ng.ProcessString(":", EnglishMode)
	ng.RemoveLastChar(true)
	if ng.GetProcessedString(VietnameseMode) != "lọa" {
		t.Errorf("Process [loanj-1], got [%s] expected [lọa]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessUpperString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("VIEETJ", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "VIỆT" {
		t.Errorf("Process [VIEETJ], got [%s] expected [VIỆT]", ng.GetProcessedString(VietnameseMode))
	}
	ng.RemoveLastChar(false)
	if ng.GetProcessedString(VietnameseMode) != "VIỆ" {
		t.Errorf("Process remove last char of upper string, got [%s] expected [VIỆ]", ng.GetProcessedString(VietnameseMode))
	}
	ng.ProcessKey('Q', VietnameseMode)
	if ng.GetProcessedString(EnglishMode) != "VIEEJQ" {
		t.Errorf("Process remove last char of upper string, got [%s] expected [VIEEJQ]", ng.GetProcessedString(EnglishMode))
	}
	ng.Reset()
	ng.ProcessString("IB", EnglishMode)
	if ng.GetProcessedString(EnglishMode) != "IB" {
		t.Errorf("Process remove last char of upper string, got [%s] expected [IB]", ng.GetProcessedString(EnglishMode))
	}
}

func TestSpellingCheck(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("noww", VietnameseMode)
	if ng.GetProcessedString(EnglishMode) != "noww" {
		t.Errorf("Process-ENG [noww], got [%s] expected [noww]", ng.GetProcessedString(EnglishMode))
	}
	if ng.GetProcessedString(VietnameseMode) != "now" {
		t.Errorf("Process-VIE [noww], got [%s] expected [now]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("sawss", VietnameseMode)
	if ng.GetProcessedString(EnglishMode) != "sawss" {
		t.Errorf("Process-ENG [sawss], got [%s] expected [sawss]", ng.GetProcessedString(EnglishMode))
	}
	ng.Reset()
	ng.ProcessString("sawss", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "săs" {
		t.Errorf("Process-VIE [sawss], got [%s] expected [săs]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessDD(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("dd", VietnameseMode)
	if ng.IsValid(false) == false {
		t.Errorf("Check spelling for [dd], got [%v] expected [true]", ng.IsValid(false) == false)
	}
	if ng.GetProcessedString(VietnameseMode) != "đ" {
		t.Errorf("Process [dd], got [%s] expected [đ]", ng.GetProcessedString(EnglishMode))
	}
	ng.Reset()
	ng.ProcessString("SD", VietnameseMode)
	ng.ProcessString("D", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "SĐ" {
		t.Errorf("IsSpellingCorrect [SDD], got [%v] expected [SĐ]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestTelex23(t *testing.T) {
	ng = newStdEngine()
	ng.ProcessString("t ]", EnglishMode)
	ng.ProcessString("a", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "]a" {
		t.Errorf("Process t ]a, got %s valid=%v expected true", ng.GetProcessedString(VietnameseMode), ng.IsValid(false))
	}
	ng.Reset()
	ng.ProcessString("]]a", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "]a" {
		t.Errorf("Process ]aa, got %s valid=%v expected true", ng.GetProcessedString(VietnameseMode), ng.IsValid(true))
	}
	var im = ParseInputMethod(InputMethodDefinitions, "Telex 2")
	var ng = NewEngine(im, EstdFlags)
	ng.ProcessString("[", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "ơ" {
		t.Errorf("Process Telex 2 [[], got [%v] expected [ươ]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("{", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "Ơ" {
		t.Errorf("Process Telex 2 [{], got [%s] expected [Ơ]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessNguwowfiString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("wowfi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "ười" {
		t.Errorf("Process [wowfi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "ười")
	}
}

func TestRemoveLastChar(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("hanhj", VietnameseMode)
	ng.RemoveLastChar(true)
	if ng.GetProcessedString(VietnameseMode) != "hạn" {
		t.Errorf("Process [hanhj], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "hạn")
	}
	ng.Reset()
}

func TestProcessCatrString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("catr", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "catr" {
		t.Errorf("Process [nguwowfi], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "catr")
	}
}

func TestProcessToowiString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("toowi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tơi" {
		t.Errorf("Process [toowi], got [%s] expected [tơi]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessAlooString(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("aloo", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "alô" {
		t.Errorf("Process [aloo], got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode), "alô")
	}
}

func TestSpellingCheckForGiw(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("giw", VietnameseMode)
	if ng.IsValid(false) == false {
		t.Errorf("Process giw, got [%v] expected [%v]", ng.IsValid(false) == false, true)
	}
}

func TestDoubleBrackets(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("[[", VietnameseMode)
	if ng.GetProcessedString(EnglishMode) != "[" {
		t.Errorf("TestDoubleBrackets, got [%v] expected [%v]", ng.GetProcessedString(EnglishMode), "[")
	}
}
func TestDoubleBracketso(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("tooss", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tôs" {
		t.Errorf("Process tooss, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("tosos", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tôs" {
		t.Errorf("Process tosos, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestDoubleW(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("ww", VietnameseMode)
	if ng.GetProcessedString(EnglishMode) != "w" {
		t.Errorf("TestDoubleW-ENG, got [%v] expected [w]", ng.GetProcessedString(EnglishMode))
	}
	if ng.GetProcessedString(VietnameseMode) != "w" {
		t.Errorf("TestDoubleW-VIE, got [%v] expected [w]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestDoubleW2(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("wiw", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "uiw" {
		t.Errorf("TestDoubleW-VIE wiw, got [%v] expected [uiw]", ng.GetProcessedString(VietnameseMode))
	}
	if ng.GetProcessedString(EnglishMode) != "wiw" {
		t.Errorf("TestDoubleW-ENG wiw, got [%v] expected [wiw]", ng.GetProcessedString(EnglishMode))
	}
}

func TestProcessDuwoi(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("duwoi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "dươi" {
		t.Errorf("Process duwoi, got [%v] expected [dươi]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessRefresh(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("reff", VietnameseMode)
	ng.ProcessString("resh", EnglishMode)
	if ng.GetProcessedString(EnglishMode) != "reffresh" {
		t.Errorf("Process-ENG [reff+resh], got [%v] expected [reffresh]", ng.GetProcessedString(EnglishMode))
	}
	if ng.GetProcessedString(VietnameseMode) != "refresh" {
		t.Errorf("Process-VIE [reff+resh], got [%v] expected [refresh]", ng.GetProcessedString(VietnameseMode))
	}
}
func TestProcessRefresh2(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("reff", VietnameseMode)
	ng.RemoveLastChar(true)
	ng.ProcessKey('f', VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "rè" {
		t.Errorf("Process reff-1+f, got [%v] expected [rè]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessDDSeq(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("oddp", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "ođp" {
		t.Errorf("Process oddp, got [%v] expected [ođp]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessGisa(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("gis", VietnameseMode)
	ng.ProcessString("a", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "giá" {
		t.Errorf("Process gisa, got [%v] expected [giá]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessKimso(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("kimso", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "kímo" {
		t.Errorf("TestProcessKimso, got [%v] expected [kímo]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestProcessTo(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("to", VietnameseMode)
	if ng.IsValid(true) == false {
		t.Errorf("Process to, got [%v] expected [true]", ng.IsValid(true) == false)
	}
}

func TestProcessToorr(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("toorr", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tôr" {
		t.Errorf("Process toorr, got [%v] expected [tôr]", ng.GetProcessedString(VietnameseMode))
	}
}

//tnó
func TestProcessTnoss(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("tnoss", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tnos" {
		t.Errorf("Process tnoss, got [%v] expected [tnos]", ng.GetProcessedString(VietnameseMode))
	}
}

//ềng
func TestProcessEenghf(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "Telex 2")
	ng := NewEngine(im, EstdFlags)
	ng.ProcessString("ddawks", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "đắk" {
		t.Errorf("Process eenghf, got [%v] expected [đắk]", ng.GetProcessedString(VietnameseMode))
	}
}

//HIEEUR
func TestProcessHIEEUR(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("tooi oo HIEEUR", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "HIỂU" {
		t.Errorf("Process [tooi oo HIEEUR], got [%v] expected [HIỂU]", ng.GetProcessedString(VietnameseMode))
	}
}

//NGUOIW
func TestProcessNGUOIW(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("NGUOIW", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "NGƯƠI" {
		t.Errorf("TestProcessToorr, got [%v] expected [NGƯƠI]", ng.GetProcessedString(VietnameseMode))
	}
}

//T{s
func TestProcessTOs(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("{s", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "Ớ" {
		t.Errorf("Process {+s, got [%v] expected [Ớ]", ng.GetProcessedString(VietnameseMode))
	}
}

//T{s
func TestProcessTo5(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "VNI")
	ng := NewEngine(im, EstdFlags)
	ng.ProcessString("o55", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "o5" {
		t.Errorf("Process [o55-VNI], got [%v] expected [o5]", ng.GetProcessedString(VietnameseMode))
	}
}

//duwongwj
func TestProcesshuoswc(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("duwongwj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "duongwj" {
		t.Errorf("Process [duwongwj], got [%v] expected [duongwj]", ng.GetProcessedString(VietnameseMode))
	}
}

//choas, bieecs, uese
func TestProcesschoas(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "Telex 2")
	ng := NewEngine(im, EstdFlags&^EstdToneStyle)
	ng.ProcessString("choas", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "choá" {
		t.Errorf("Process [choas], got [%v] expected [choá]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("bieecs", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "biếc" {
		t.Errorf("Process [bieecs], got [%v] expected [biếc]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("uese", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "uế" {
		t.Errorf("Process uese, got [%v] expected [uế]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestBambooEngine_RestoreLastWord(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("duwongj tooi", VietnameseMode)
	ng.RestoreLastWord(false)
	if ng.GetProcessedString(VietnameseMode) != "tooi" {
		t.Errorf("Process [duwongwj tooi], got [%v] expected [tooi]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestBambooEngine_RestoreLastWord_TCVN(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "Microsoft layout")
	ng := NewEngine(im, EstdFlags)
	ng.ProcessString("112", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "1â" {
		t.Errorf("Process-VIE 112 (Microsoft layout), got [%v] expected [1â]", ng.GetProcessedString(VietnameseMode))
	}
	ng.RestoreLastWord(false)
	if ng.GetProcessedString(EnglishMode) != "12" {
		t.Errorf("Process-ENG 112 (Microsoft layout), got [%v] expected [12]", ng.GetProcessedString(EnglishMode))
	}
	ng.Reset()
	ng.ProcessString("d[]ng9 t4i", VietnameseMode)
	ng.RestoreLastWord(false)
	if ng.GetProcessedString(VietnameseMode) != "t4i" {
		t.Errorf("Process [duongwj t4i - MS layout], got [%v] expected [t4i]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestBambooEngine_Zprocessing(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("loz", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "loz" {
		t.Errorf("Process loz, got [%v] expected [loz]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("losz", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "lo" {
		t.Errorf("Process-VIE losz, got [%v] expected [lo]", ng.GetProcessedString(VietnameseMode))
	}
	if ng.GetProcessedString(EnglishMode) != "losz" {
		t.Errorf("Process-ENG losz, got [%v] expected [losz]", ng.GetProcessedString(EnglishMode))
	}
}

func TestRestoreLastWord(t *testing.T) {
	ng := newStdEngine()
	s := "afq"
	ng.ProcessString(s, VietnameseMode)
	ng.RestoreLastWord(false)
	ng.RemoveLastChar(true)
	ng.ProcessKey('f', VietnameseMode)
	t.Logf("LOGGING Process [%s] got [%v], en=[%s]", s, ng.GetProcessedString(VietnameseMode), ng.GetProcessedString(EnglishMode))
}

func TestProcessVNWord(t *testing.T) {
	var s = "tôifs"
	ng := newStdEngine()
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tối" {
		t.Errorf("Process tôifs, got [%v] expected [tối]", ng.GetProcessedString(VietnameseMode))
	}
	if ng.GetProcessedString(EnglishMode) != "tôifs" {
		t.Errorf("Process-ENG tôifs, got [%v] expected [tôifs]", ng.GetProcessedString(EnglishMode))
	}
	ng.Reset()
	s = "tốif"
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tồi" {
		t.Errorf("Process tôifs, got [%v] expected [tồi]", ng.GetProcessedString(VietnameseMode))
	}
	if ng.GetProcessedString(EnglishMode) != "tốif" {
		t.Errorf("Process tôifs, got [%v] expected [tốif]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	s = "tốiz"
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tôi" {
		t.Errorf("Process tôifs, got [%v] expected [tôi]", ng.GetProcessedString(VietnameseMode))
	}
}

func TestDoubleTyping(t *testing.T) {
	var s = "linux"
	ng := newStdEngine()
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("x", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "linux" {
		t.Errorf("Process [linuxx], got [%v] expected [linux]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	s = "buwo"
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("o", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "buô" {
		t.Errorf("Process [buwoo], got [%s] expected [buô]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	s = "buowc"
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("o", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "buôc" {
		t.Errorf("Process [buowco], got [%s] expected [buôc]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	s = "cuoiw"
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("o", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "cuôi" {
		t.Errorf("Process [cuoiwo], got [%s] expected [cuôi]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	s = "ach"
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("a", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "acha" {
		t.Errorf("Process [acha], got [%s] expected [acha]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	s = "nhuw"
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "như" {
		t.Errorf("Process [acha], got [%s] expected [như]", ng.GetProcessedString(VietnameseMode))
	}
	if ng.IsValid(true) != true {
		t.Errorf("Findresultmatch full, got %v expected true", ng.IsValid(true))
	}
	// AddDictionaryToSpellingTrie(map[string]bool{"thứ": true})
	ng.Reset()
	s = "thuw"
	ng.ProcessString(s, VietnameseMode)
	if ng.IsValid(true) != true {
		t.Errorf("true, got %v expected true", ng.IsValid(true))
	}
	ng.Reset()
	s = "thow"
	ng.ProcessString(s, VietnameseMode)
	if ng.IsValid(true) != true {
		t.Errorf("true, got %v expected true", ng.IsValid(true))
	}
	ng.Reset()
	// AddDictionaryToSpellingTrie(map[string]bool{"tôi": true, "tối": true, "tời": true, "tơi": true})
	s = "tooi"
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tôi" {
		t.Errorf("Process [acha], got [%s] expected [tôi]", ng.GetProcessedString(VietnameseMode))
	}
	if ng.IsValid(true) != true {
		t.Errorf("Findresultmatch full, got %v expected true", ng.IsValid(true))
	}
	ng.Reset()
	ng.ProcessString("arch", VietnameseMode)
	if ng.IsValid(false) != false {
		t.Errorf("false arch, got %v expected 0", ng.IsValid(false))
	}
	ng.Reset()
	ng.ProcessString("[[", VietnameseMode)
	ng.ProcessString("oo", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "[ô" {
		t.Errorf("Process [oo, got %s expected [ô", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("oo]", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "ôư" {
		t.Errorf("Process oo], got %s expected ôư", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("chury", VietnameseMode)
	if ng.IsValid(true) == false {
		t.Errorf("IsValid chury, got %v expected 0", ng.IsValid(true))
	}
	ng.Reset()
	ng.ProcessString("turyn", VietnameseMode)
	ng.RemoveLastChar(true)
	ng.RemoveLastChar(true)
	// ng.ProcessString("r", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "tủ" {
		t.Errorf("Process turyen,BS,BS,BS,r, got [%s] expected [tủ]", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("chuyển", VietnameseMode)
	ng.ProcessString("z", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "chuyên" {
		t.Errorf("Process [chuyểnz], got %s expected chuyên", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("nhueej", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "nhuệ" {
		t.Errorf("Process nhueej, got %s expected nhuệ", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("cuongw", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "cương" {
		t.Errorf("Process cuongw, got %s expected cương", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
	ng.ProcessString("quawcj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "quặc" {
		t.Errorf("Process quawcj, got %s expected quặc", ng.GetProcessedString(VietnameseMode))
	}
	//eechs
	ng.Reset()
	ng.ProcessString("quawcj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode) != "quặc" {
		t.Errorf("Process quawcj, got %s valid=%v expected quặc", ng.GetProcessedString(VietnameseMode), ng.IsValid(false))
	}
	ng.Reset()
	ng.ProcessString("tôi）t", EnglishMode)
	if ng.GetProcessedString(VietnameseMode) != "t" {
		t.Errorf("Process [tôi）t], got %s expected t", ng.GetProcessedString(VietnameseMode))
	}
	ng.Reset()
}

var ng = newStdEngine()

func BenchmarkRemoveLastChar(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	ng.Reset()
	for i := 0; i < b.N; i++ {
		ng.ProcessString(" ", EnglishMode)
		ng.ProcessString("aj", VietnameseMode)
		if ng.GetProcessedString(VietnameseMode) != "ạ" {
			b.Errorf("Process [aj], got [%s] expected [ạ]", ng.GetProcessedString(VietnameseMode))
		}
	}
}
