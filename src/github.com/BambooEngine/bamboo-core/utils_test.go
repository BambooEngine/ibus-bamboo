package bamboo

import (
	"testing"
)

func TestIsVowel(t *testing.T) {
	if IsVowel('a') == false {
		t.Errorf("a is a vowel, but result is false")
	}
	if IsVowel('á') == false {
		t.Errorf("á is a vowel, but result is false")
	}
	if IsVowel('b') {
		t.Errorf("b is not a vowel, but result is true")
	}
	tvowels := []rune("aàáảãạăằắẳẵặâầấẩẫậeèéẻẽẹêềếểễệiìíỉĩịoòóỏõọôồốổỗộơờớởỡợuùúủũụưừứửữựyỳýỷỹỵ")
	for _, v := range tvowels {
		if IsVowel(v) == false {
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
	c_y := AddToneToChar('y', 0)
	if c_y != 'y' {
		t.Errorf("Add TONE_NONE to char y, got %c expected y", c_y)
	}
	c_y = AddMarkToChar('y', 0)
	if c_y != 'y' {
		t.Errorf("Add MARK_NONE to char y, got %c expected y", c_y)
	}
}

func TestAddMarkToChar(t *testing.T) {
	c_ặ := AddMarkToChar('ạ', uint8(MARK_BREVE))
	if c_ặ != 'ặ' {
		t.Errorf("Test add a breve to char. Got %c, expected %c", c_ặ, 'ặ')
	}
}
