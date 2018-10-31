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

var MARKS_MAP = map[rune]string{
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

func RemoveMarkFromChar(chr rune) rune {
	if str, found := MARKS_MAP[chr]; found {
		marks := []rune(str)
		if len(marks) > 0 {
			return marks[0]
		}
	}
	return chr
}

func AddMarkToChar(chr rune, mark Mark) rune {
	var result rune
	tone := FindToneFromChar(chr)
	chr = AddToneToChar(chr, TONE_NONE)
	if str, found := MARKS_MAP[chr]; found {
		marks := []rune(str)
		if marks[mark] != '_' {
			result = marks[mark]
		}
	}
	result = AddToneToChar(result, tone)
	return result
}

func FindMarkTargets(composition []*Transformation, rule Rule) []*Transformation {
	var result []*Transformation
	for _, trans := range composition {
		if trans.Rule.Key == rule.EffectOn {
			result = append(result, trans)
		}
	}
	return result
}

func FindMarkTarget(composition []*Transformation, rules []Rule) (*Transformation, Rule) {
	for i := len(composition) - 1; i >= 0; i-- {
		var trans = composition[i]
		for _, rule := range rules {
			if trans.Rule.Key == rule.EffectOn {
				var target = trans
				if isMarkTargetValid(composition, &Transformation{
					Rule: rule, Target: target}) {
					return target, rule
				}
			}
		}
	}
	return nil, Rule{}
}

func isMarkTargetValid(composition []*Transformation, trans *Transformation) bool {
	if !isFree(composition, trans.Target, MarkTransformation) {
		return false
	}
	var soundMap = GetSoundMap(composition)
	targetSound, found := soundMap[trans.Target]
	if !found {
		return false
	}
	if isVowel(trans.Rule.EffectOn) && targetSound != VowelSound {
		return false
	}
	if targetSound == VowelSound && !isSpellingCorrect(GetRightMostVowelWithMarks(append(composition, trans)), NoTone) {
		return false
	}
	return true
}

func GetMarkTransformationsTargetTo(composition []*Transformation, trans *Transformation) []*Transformation {
	var result []*Transformation
	for _, t := range composition {
		if t.Target == trans && t.Rule.EffectType == MarkTransformation {
			result = append(result, t)
		}
	}
	return result
}

func GetTransformationsTargetTo(composition []*Transformation, trans *Transformation) []*Transformation {
	var result []*Transformation
	for _, t := range composition {
		if t.Target == trans {
			result = append(result, t)
		}
	}
	return result
}

func GetRightMostVowelWithMarks(composition []*Transformation) []*Transformation {
	var vowels = GetRightMostVowels(composition)
	return AddMarksToComposition(composition, vowels)
}

func AddMarksToComposition(composition []*Transformation, appendingComps []*Transformation) []*Transformation {
	var result []*Transformation
	result = append(result, appendingComps...)
	for _, t := range appendingComps {
		result = append(result, GetMarkTransformationsTargetTo(composition, t)...)
	}
	return result
}

func GetLastCombinationWithMarks(composition []*Transformation) []*Transformation {
	return AddMarksToComposition(composition, GetLastCombination(composition))
}

var vowelMap = map[string]rune{
	"uo": 'o',
	"ua": 'u',
	"ou": 'o',
	"oa": 'a',
	"au": 'a',
	"ao": 'a',
}

func FindMarkTargetIndex(chars []rune) int {
	if len(chars) != 2 {
		return 0
	}
	if chars[0] == chars[1] {
		return 0
	}
	if key, found := vowelMap[string(chars)]; found {
		if key == chars[0] {
			return 0
		}
		return 1
	}
	return 0
}
