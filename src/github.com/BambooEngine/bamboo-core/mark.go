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
	if str, found := marksMaps[chr]; found {
		marks := []rune(str)
		if marks[mark] != '_' {
			result = marks[mark]
		}
	}
	result = AddToneToChar(result, uint8(tone))
	return result
}

func findMarkTargets(composition []*Transformation, rule Rule) []*Transformation {
	var result []*Transformation
	for _, trans := range composition {
		if trans.Rule.Key == rule.EffectOn {
			result = append(result, trans)
		}
	}
	return result
}

func findMarkTarget(composition []*Transformation, rules []Rule) (*Transformation, Rule) {
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
	// the target of trans is a vowel
	if IsVowel(trans.Rule.EffectOn) && targetSound != VowelSound {
		return false
	}
	if targetSound == VowelSound && !isSpellingCorrect(getRightMostVowelWithMarks(append(composition, trans)), NoTone) {
		return false
	}
	return true
}

func getMarkTransformationsTargetTo(composition []*Transformation, trans *Transformation) []*Transformation {
	var result []*Transformation
	for _, t := range composition {
		if t.Target == trans && t.Rule.EffectType == MarkTransformation {
			result = append(result, t)
		}
	}
	return result
}

func getTransformationsTargetTo(composition []*Transformation, trans *Transformation) []*Transformation {
	var result []*Transformation
	for _, t := range composition {
		if t.Target == trans {
			result = append(result, t)
		}
	}
	return result
}

func addMarksToComposition(composition []*Transformation, appendingComps []*Transformation) []*Transformation {
	var result []*Transformation
	result = append(result, appendingComps...)
	for _, t := range appendingComps {
		result = append(result, getMarkTransformationsTargetTo(composition, t)...)
	}
	return result
}
