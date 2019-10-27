/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENCE>.
 */
package bamboo

import "unicode"

var Vowels = []rune("aàáảãạăằắẳẵặâầấẩẫậeèéẻẽẹêềếểễệiìíỉĩịoòóỏõọôồốổỗộơờớởỡợuùúủũụưừứửữựyỳýỷỹỵ")

var PunctuationMarks = []rune{
	',', ';', ':', '.', '"', '\'', '!', '?', ' ',
	'<', '>', '=', '+', '-', '*', '/', '\\',
	'_', '~', '`', '@', '#', '$', '%', '^', '&', '(', ')', '{', '}', '[', ']',
	'|',
}

func IsPunctuationMark(key rune) bool {
	for _, c := range PunctuationMarks {
		if c == key {
			return true
		}
	}
	return false
}

func IsWordBreakSymbol(key rune) bool {
	return IsPunctuationMark(key) || ('0' <= key && '9' >= key)
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

func FindVowelPosition(chr rune) int {
	for pos, v := range Vowels {
		if v == chr {
			return pos
		}
	}
	return -1
}

var marksMaps = map[rune]string{
	'a': "aâă__",
	'â': "aâă__",
	'ă': "aâă__",
	'e': "eê___",
	'ê': "eê___",
	'o': "oô_ơ_",
	'ô': "oô_ơ_",
	'ơ': "oô_ơ_",
	'u': "u__ư_",
	'ư': "u__ư_",
	'd': "d___đ",
	'đ': "d___đ",
}

func getMarkFamily(chr rune) []rune {
	var result []rune
	if s, found := marksMaps[chr]; found {
		for _, c := range []rune(s) {
			if c != '_' {
				result = append(result, c)
			}
		}
	}
	return result
}

func FindMarkPosition(chr rune) int {
	if str, found := marksMaps[chr]; found {
		for pos, v := range []rune(str) {
			if v == chr {
				return pos
			}
		}
	}
	return -1
}

func FindMarkFromChar(chr rune) (Mark, bool) {
	var pos = FindMarkPosition(chr)
	if pos >= 0 {
		return Mark(pos), true
	}
	return 0, false
}

func AddMarkToTonelessChar(chr rune, mark uint8) rune {
	if str, found := marksMaps[chr]; found {
		marks := []rune(str)
		if marks[mark] != '_' {
			return marks[mark]
		}
	}
	return chr
}

func AddMarkToChar(chr rune, mark uint8) rune {
	tone := FindToneFromChar(chr)
	chr = AddToneToChar(chr, 0)
	chr = AddMarkToTonelessChar(chr, mark)
	return AddToneToChar(chr, uint8(tone))
}

func IsAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func inKeyList(keys []rune, key rune) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func FindToneFromChar(chr rune) Tone {
	pos := FindVowelPosition(chr)
	if pos == -1 {
		return TONE_NONE
	}
	return Tone(pos % 6)
}

func AddToneToChar(chr rune, tone uint8) rune {
	pos := FindVowelPosition(chr)
	if pos > -1 {
		current_tone := pos % 6
		offset := int(tone) - current_tone
		return Vowels[pos+offset]
	} else {
		return chr
	}
}

func IsVietnameseRune(chr rune) bool {
	var c = unicode.ToLower(chr)
	if FindToneFromChar(c) != TONE_NONE {
		return true
	}
	return c != AddMarkToTonelessChar(c, 0)
}

func HasAnyVietnameseRune(word string) bool {
	for _, chr := range []rune(word) {
		if IsVietnameseRune(chr) {
			return true
		}
	}
	return false
}
