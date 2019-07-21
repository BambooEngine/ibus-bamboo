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
	if ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch {
		t.Errorf("IsSpellingCorrect [dd], got [%v] expected [true]", ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch)
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
		t.Errorf("Process [Thuow] and remove last char, got [%s] expected [%s]", ng.GetProcessedString(VietnameseMode, false), "Thu")
	}
}

func TestBambooEngine_RemoveLastChar(t *testing.T) {
	ng := newStdEngine()
	ng.RemoveLastChar()
	ng.ProcessString(" ", EnglishMode)
	ng.RemoveLastChar()
	ng.ProcessString("loanj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "loạn" {
		t.Errorf("Process [loanj], got [%s] expected [loạn]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode, false) != "lọa" {
		t.Errorf("Process [loanj-1], got [%s] expected [lọa]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.ProcessString(":", EnglishMode)
	ng.RemoveLastChar()
	if ng.GetProcessedString(VietnameseMode, false) != "lọa" {
		t.Errorf("Process [loanj-1], got [%s] expected [lọa]", ng.GetProcessedString(VietnameseMode, false))
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
	ng.ProcessKey('Q', VietnameseMode)
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
	if ng.GetProcessedString(EnglishMode, false) != "noww" {
		t.Errorf("Process-ENG [noww], got [%s] expected [noww]", ng.GetProcessedString(EnglishMode, false))
	}
	if ng.GetProcessedString(VietnameseMode, false) != "now" {
		t.Errorf("Process-VIE [noww], got [%s] expected [now]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	ng.ProcessString("sawss", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "sawss" {
		t.Errorf("Process-ENG [sawss], got [%s] expected [sawss]", ng.GetProcessedString(EnglishMode, false))
	}
	ng.Reset()
	ng.ProcessString("sawss", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "săs" {
		t.Errorf("Process-VIE [sawss], got [%s] expected [săs]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessDD(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("dd", VietnameseMode)
	if ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch {
		t.Errorf("Check spelling for [dd], got [%v] expected [true]", ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch)
	}
	if ng.GetProcessedString(VietnameseMode, false) != "đ" {
		t.Errorf("Process [dd], got [%s] expected [đ]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestTelex3(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "Telex 3")
	var ng = NewEngine(im, EstdFlags)
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
	if ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch {
		t.Errorf("Process giw, got [%v] expected [%v]", ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch, true)
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
		t.Errorf("Process tooss, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	ng.ProcessString("tosos", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tôs" {
		t.Errorf("Process tosos, got [%v] expected [tôs]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestDoubleW(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("ww", VietnameseMode)
	if ng.GetProcessedString(EnglishMode, false) != "w" {
		t.Errorf("TestDoubleW-ENG, got [%v] expected [w]", ng.GetProcessedString(EnglishMode, false))
	}
	if ng.GetProcessedString(VietnameseMode, false) != "w" {
		t.Errorf("TestDoubleW-VIE, got [%v] expected [w]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestDoubleW2(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("wiw", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "uiw" {
		t.Errorf("TestDoubleW-VIE wiw, got [%v] expected [uiw]", ng.GetProcessedString(VietnameseMode, false))
	}
	if ng.GetProcessedString(EnglishMode, false) != "wiw" {
		t.Errorf("TestDoubleW-ENG wiw, got [%v] expected [wiw]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestProcessDuwoi(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("duwoi", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "dươi" {
		t.Errorf("Process duwoi, got [%v] expected [dươi]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessRefresh(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("reff", VietnameseMode)
	ng.ProcessString("resh", EnglishMode)
	if ng.GetProcessedString(EnglishMode, false) != "reffresh" {
		t.Errorf("Process-ENG [reff+resh], got [%v] expected [reffresh]", ng.GetProcessedString(EnglishMode, false))
	}
	if ng.GetProcessedString(VietnameseMode, false) != "refresh" {
		t.Errorf("Process-VIE [reff+resh], got [%v] expected [refresh]", ng.GetProcessedString(VietnameseMode, false))
	}
}
func TestProcessRefresh2(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("reff", VietnameseMode)
	ng.RemoveLastChar()
	ng.ProcessKey('f', VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "rè" {
		t.Errorf("Process reff-1+f, got [%v] expected [rè]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessDDSeq(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("oddp", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ođp" {
		t.Errorf("Process oddp, got [%v] expected [ođp]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestProcessGisa(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("gisa", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "giá" {
		t.Errorf("Process gisa, got [%v] expected [giá]", ng.GetProcessedString(VietnameseMode, false))
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
	if ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch {
		t.Errorf("Process to, got [%v] expected [true]", ng.GetSpellingMatchResult(ToneLess, false) == FindResultNotMatch)
	}
}

func TestProcessToorr(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("toorr", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tôr" {
		t.Errorf("Process toorr, got [%v] expected [tôr]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//tnó
func TestProcessTnoss(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("tnoss", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tnos" {
		t.Errorf("Process tnoss, got [%v] expected [tnos]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//ềng
func TestProcessEenghf(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "Telex 2")
	ng := NewEngine(im, EstdFlags)
	ng.AddDictionary(map[string]bool{"ềngh": true})
	ng.ProcessString("eenghf", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "ềngh" {
		t.Errorf("Process eenghf, got [%v] expected [ềnhg]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//HIEEUR
func TestProcessHIEEUR(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("tooi oo HIEEUR", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "HIỂU" {
		t.Errorf("Process [tooi oo HIEEUR], got [%v] expected [HIỂU]", ng.GetProcessedString(VietnameseMode, false))
	}
	if ng.GetProcessedString(VietnameseMode, true) != "HIỂU" {
		t.Errorf("Process [tooi oo HIEUR], got [%v] expected [HIỂU]", ng.GetProcessedString(VietnameseMode, true))
	}
}

//NGUOIW
func TestProcessNGUOIW(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("NGUOIW", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "NGƯƠI" {
		t.Errorf("TestProcessToorr, got [%v] expected [NGƯƠI]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//T{s
func TestProcessTOs(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("{s", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "Ớ" {
		t.Errorf("Process {+s, got [%v] expected [Ớ]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//T{s
func TestProcessTo5(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "VNI")
	ng := NewEngine(im, EstdFlags)
	ng.ProcessString("o55", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "o5" {
		t.Errorf("Process [o55-VNI], got [%v] expected [o5]", ng.GetProcessedString(VietnameseMode, false))
	}
	if ng.GetProcessedString(VietnameseMode, true) != "" {
		t.Errorf("Process [o55-VNI], got [%v] expected []", ng.GetProcessedString(VietnameseMode, true))
	}
}

//duwongwj
func TestProcesshuoswc(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("duwongwj", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "duongwj" {
		t.Errorf("Process [duwongwj], got [%v] expected [duongwj]", ng.GetProcessedString(VietnameseMode, false))
	}
}

//choas, bieecs, uese
func TestProcesschoas(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "Telex 2")
	ng := NewEngine(im, EstdFlags&^EstdToneStyle)
	ng.ProcessString("choas", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "choá" {
		t.Errorf("Process [choas], got [%v] expected [choá]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	ng.ProcessString("bieecs", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "biếc" {
		t.Errorf("Process [bieecs], got [%v] expected [biếc]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	ng.ProcessString("uese", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "uế" {
		t.Errorf("Process uese, got [%v] expected [uế]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestBambooEngine_RestoreLastWord(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("duwongwj tooi", VietnameseMode)
	ng.RestoreLastWord()
	if ng.GetProcessedString(VietnameseMode, false) != "tooi" {
		t.Errorf("Process [duwongwj tooi], got [%v] expected [tooi]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestBambooEngine_RestoreLastWord_TCVN(t *testing.T) {
	var im = ParseInputMethod(InputMethodDefinitions, "Microsoft layout")
	ng := NewEngine(im, EstdFlags)
	ng.ProcessString("112", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "1â" {
		t.Errorf("Process-VIE 112 (Microsoft layout), got [%v] expected [1â]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.RestoreLastWord()
	if ng.GetProcessedString(EnglishMode, false) != "12" {
		t.Errorf("Process-ENG 112 (Microsoft layout), got [%v] expected [12]", ng.GetProcessedString(EnglishMode, false))
	}
	ng.Reset()
	ng.ProcessString("duwongwj t4i", VietnameseMode)
	ng.RestoreLastWord()
	if ng.GetProcessedString(VietnameseMode, false) != "t4i" {
		t.Errorf("Process [duwongwj t4i - MS layout], got [%v] expected [t4i]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestBambooEngine_Zprocessing(t *testing.T) {
	ng := newStdEngine()
	ng.ProcessString("loz", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "loz" {
		t.Errorf("Process loz, got [%v] expected [loz]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	ng.ProcessString("losz", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "lo" {
		t.Errorf("Process-VIE losz, got [%v] expected [lo]", ng.GetProcessedString(VietnameseMode, false))
	}
	if ng.GetProcessedString(EnglishMode, false) != "losz" {
		t.Errorf("Process-ENG losz, got [%v] expected [losz]", ng.GetProcessedString(EnglishMode, false))
	}
}

func TestRestoreLastWord(t *testing.T) {
	ng := newStdEngine()
	s := "afq"
	ng.ProcessString(s, VietnameseMode)
	ng.RestoreLastWord()
	ng.RemoveLastChar()
	ng.ProcessKey('f', VietnameseMode)
	t.Logf("Process [%s] got [%v], en=[%s]", s, ng.GetProcessedString(VietnameseMode, false), ng.GetProcessedString(EnglishMode, false))
}

func TestProcessVNWord(t *testing.T) {
	var s = "tôifs"
	ng := newStdEngine()
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tối" {
		t.Errorf("Process tôifs, got [%v] expected [tối]", ng.GetProcessedString(VietnameseMode, false))
	}
	if ng.GetProcessedString(EnglishMode, false) != "tôifs" {
		t.Errorf("Process-ENG tôifs, got [%v] expected [tôifs]", ng.GetProcessedString(EnglishMode, false))
	}
	ng.Reset()
	s = "tốif"
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tồi" {
		t.Errorf("Process tôifs, got [%v] expected [tồi]", ng.GetProcessedString(VietnameseMode, false))
	}
	if ng.GetProcessedString(EnglishMode, false) != "tốif" {
		t.Errorf("Process tôifs, got [%v] expected [tốif]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	s = "tốiz"
	ng.ProcessString(s, VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "tôi" {
		t.Errorf("Process tôifs, got [%v] expected [tôi]", ng.GetProcessedString(VietnameseMode, false))
	}
}

func TestDoubleTyping(t *testing.T) {
	var s = "linux"
	ng := newStdEngine()
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("x", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "linux" {
		t.Errorf("Process [linuxx], got [%v] expected [linux]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	s = "buowc"
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("o", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "buôc" {
		t.Errorf("Process [buowco], got [%s] expected [buôc]", ng.GetProcessedString(VietnameseMode, false))
	}
	ng.Reset()
	s = "cuoiw"
	ng.ProcessString(s, VietnameseMode)
	ng.ProcessString("o", VietnameseMode)
	if ng.GetProcessedString(VietnameseMode, false) != "cuôi" {
		t.Errorf("Process [cuoiw], got [%s] expected [cuôi]", ng.GetProcessedString(VietnameseMode, false))
	}
}
