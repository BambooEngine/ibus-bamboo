package bamboo

import "strings"

var vowels = []rune("aàáảãạăằắẳẵặâầấẩẫậeèéẻẽẹêềếểễệiìíỉĩịoòóỏõọôồốổỗộơờớởỡợuùúủũụưừứửữựyỳýỷỹỵ")

var vowelMap = map[string]rune{
	"uo": 'o',
	"ua": 'u',
	"ou": 'o',
	"oa": 'a',
	"au": 'a',
	"ao": 'a',
}

func isVowel(chr rune) bool {
	isVowel := false
	for _, v := range vowels {
		if v == chr {
			isVowel = true
		}
	}
	return isVowel
}

func FindVowelPosition(chr rune) int {
	for pos, v := range vowels {
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

func isVowelString(str string) bool {
	var isVowels = true
	for _, chr := range []rune(str) {
		if !isVowel(chr) {
			isVowels = false
		}
	}
	return isVowels
}

func findVowelByOrder(vowels []*Transformation, order int) *Transformation {
	var i = 0
	for _, trans := range vowels {
		if trans.Rule.EffectType == Appending {
			if i == order {
				return trans
			}
			i++
		}
	}
	return nil
}

func getRightMostVowels(composition []*Transformation) []*Transformation {
	return GetLastSoundGroup(composition, VowelSound)
}
