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

func FindToneFromChar(chr rune) Tone {
	pos := FindVowelPosition(chr)
	if pos == -1 {
		return TONE_NONE
	}
	return Tone(pos % 6)
}

func AddToneToChar(chr rune, tone Tone) rune {
	pos := FindVowelPosition(chr)
	if pos > -1 {
		current_tone := pos % 6
		offset := int(tone) - current_tone
		return vowelSeq[pos+offset]
	} else {
		return chr
	}
}

func FindToneTarget(composition []*Transformation, stdStyle bool) *Transformation {
	if len(composition) == 0 {
		return nil
	}
	var target *Transformation
	var vowels = GetRightMostVowels(composition)
	if len(vowels) == 1 {
		target = vowels[0]
	} else if len(vowels) == 2 && stdStyle {
		var str = Flatten(GetRightMostVowelWithMarks(composition), NoTone|LowerCase)
		var chars = []rune(str)
		if ohPos := FindIndexRune(chars, 'ơ'); ohPos > 0 {
			target = vowels[ohPos]
		} else if ehPos := FindIndexRune(chars, 'ê'); ehPos > 0 {
			target = vowels[ehPos]
		} else {
			if _, found := FindNextAppendingTransformation(composition, vowels[1]); found {
				target = vowels[1]
			} else {
				target = vowels[0]
			}
		}
	} else if len(vowels) == 2 {
		var str = Flatten(GetRightMostVowels(composition), EnglishMode|LowerCase)
		if str == "oa" || str == "oe" || str == "uy" {
			target = vowels[1]
		} else {
			target = vowels[0]
		}
	} else if len(vowels) == 3 {
		if Flatten(vowels, EnglishMode|LowerCase) == "uye" {
			target = vowels[2]
		} else {
			target = vowels[1]
		}
	}
	return target
}

func hasValidTone(composition []*Transformation, tone Tone) bool {
	if tone == TONE_ACUTE || tone == TONE_DOT {
		return true
	}
	var lastConsonants = Flatten(GetLastSoundGroup(composition, LastConsonantSound), EnglishMode|LowerCase)

	// These consonants can only go with ACUTE, DOT accents
	var dotWithConsonants = []string{"c", "p", "t", "ch"}
	for _, s := range dotWithConsonants {
		if s == lastConsonants {
			return false
		}
	}
	return true
}

func GetLastToneTransformation(composition []*Transformation) *Transformation {
	for i := len(composition) - 1; i >= 0; i-- {
		var t = composition[i]
		if t.Rule.EffectType == ToneTransformation && t.Target != nil {
			return t
		}
	}
	return nil
}

func shouldRefreshLastToneTarget(transformations []*Transformation) bool {
	var vowels = GetRightMostVowels(transformations)
	if len(vowels) <= 0 {
		return false
	}
	return len(transformations) > 0 && transformations[len(transformations)-1].Rule.EffectType == Appending
}
