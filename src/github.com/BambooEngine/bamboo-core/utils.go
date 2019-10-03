/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LISENCE>.
 */
package bamboo

import "unicode"

var Vowels = []rune("aàáảãạăằắẳẵặâầấẩẫậeèéẻẽẹêềếểễệiìíỉĩịoòóỏõọôồốổỗộơờớởỡợuùúủũụưừứửữựyỳýỷỹỵ")

var WordBreakSymbols = []rune{
	',', ';', ':', '.', '"', '\'', '!', '?', ' ',
	'<', '>', '=', '+', '-', '*', '/', '\\',
	'_', '~', '`', '@', '#', '$', '%', '^', '&', '(', ')', '{', '}', '[', ']',
	'|',
}

func IsWordBreakSymbol(key rune) bool {
	for _, c := range WordBreakSymbols {
		if c == key {
			return true
		}
	}
	return false
}

func ContainsWordBreakSymbol(s string) bool {
	for _, c := range s {
		if IsWordBreakSymbol(c) {
			return true
		}
	}
	return false
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

func HasAnyVowel(seq []rune) bool {
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

func RemoveMarkFromChar(chr rune) rune {
	if str, found := marksMaps[chr]; found {
		marks := []rune(str)
		if len(marks) > 0 {
			return marks[0]
		}
	}
	return chr
}

func AddMarkToChar(chr rune, mark uint8) rune {
	var result rune
	tone := FindToneFromChar(chr)
	chr = AddToneToChar(chr, 0)
	result = chr
	if str, found := marksMaps[chr]; found {
		marks := []rune(str)
		if marks[mark] != '_' {
			result = marks[mark]
		}
	}
	result = AddToneToChar(result, uint8(tone))
	return result
}

func IsAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func findIndexRune(chars []rune, r rune) int {
	for i, c := range chars {
		if c == r {
			return i
		}
	}
	return -1
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
	if mark, found := FindMarkFromChar(AddToneToChar(c, 0)); found && mark != MARK_NONE {
		return true
	}
	return false
}

func HasAnyVietnameseRune(word string) bool {
	for _, chr := range []rune(word) {
		if IsVietnameseRune(chr) {
			return true
		}
	}
	return false
}
