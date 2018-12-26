package bamboo

import "strings"

var Vowels = []rune("aàáảãạăằắẳẵặâầấẩẫậeèéẻẽẹêềếểễệiìíỉĩịoòóỏõọôồốổỗộơờớởỡợuùúủũụưừứửữựyỳýỷỹỵ")

var vowelMap = map[string]rune{
	"uo": 'o',
	"ua": 'u',
	"ou": 'o',
	"oa": 'a',
	"au": 'a',
	"ao": 'a',
}

func IsVowel(chr rune) bool {
	isVowel := false
	for _, v := range Vowels {
		if v == chr {
			isVowel = true
		}
	}
	return isVowel
}

func HasVowel(seq []rune) bool {
	for _, s := range seq {
		if IsVowel(s) {
			return true
		}
	}
	return false
}

func FindVowelPosition(chr rune) int {
	for pos, v := range Vowels {
		if v == chr {
			return pos
		}
	}
	return -1
}

func isVowelSound(str string) bool {
	for _, line := range vowelSeq {
		for _, vowel := range strings.Split(line, " ") {
			if str == vowel {
				return true
			}
		}
	}
	return false
}

func IsVowelString(str string) bool {
	var isVowels = true
	for _, chr := range []rune(str) {
		if !IsVowel(chr) {
			isVowels = false
		}
	}
	return isVowels
}
