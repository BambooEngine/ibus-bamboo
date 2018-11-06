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

import (
	"log"
	"strings"
)

func FilterComposition(composition []*Transformation, effectType EffectType) []*Transformation {
	var result []*Transformation
	for _, trans := range composition {
		if trans.Rule.EffectType == effectType {
			result = append(result, trans)
		}
	}
	return result
}

func SeparateComposition(composition []*Transformation) [][]*Transformation {
	var result [][]*Transformation
	var seq []*Transformation
	var appendingTransformations = FilterComposition(composition, Appending)
	for i, trans := range appendingTransformations {
		seq = append(seq, trans)
		if i+1 < len(seq)-1 && isVowel(trans.Rule.EffectOn) && !isVowel(seq[i+1].Rule.EffectOn) {
			result = append(result, seq)
			seq = []*Transformation{}
		}
	}
	if len(seq) > 0 {
		result = append(result, seq)
	}
	return result
}

func belongToComposition(composition []*Transformation, trans *Transformation) bool {
	for _, t := range composition {
		if t == trans {
			return true
		}
	}
	return false
}
func GetAppendingComposition(composition []*Transformation) []*Transformation {
	var appendingTransformations []*Transformation
	for _, trans := range composition {
		if trans.Rule.EffectType == Appending {
			appendingTransformations = append(appendingTransformations, trans)
		}
	}
	return appendingTransformations
}

func isFree(composition []*Transformation, trans *Transformation, effectType EffectType) bool {
	for _, t := range composition {
		if t.Target == trans && t.Rule.EffectType == effectType {
			return false
		}
	}
	return true
}

func GetSoundMap(composition []*Transformation) map[*Transformation]Sound {
	var soundMap = map[*Transformation]Sound{}
	var appendingTransformations = GetAppendingComposition(composition)
	var str = Flatten(composition, NoTone|LowerCase)
	var sounds = ParseSoundsFromTonelessWord(str)
	if len(sounds) <= 0 || len(sounds) != len(appendingTransformations) {
		log.Println("Something wrong with the length of sounds")
		return soundMap
	}
	for i, trans := range appendingTransformations {
		soundMap[trans] = sounds[i]
	}
	return soundMap
}

func FindTransPos(composition []*Transformation, trans *Transformation) int {
	for i, t := range composition {
		if t == trans {
			return i
		}
	}
	return -1
}

func GetLastCombination(composition []*Transformation) []*Transformation {
	var combinations = SeparateComposition(composition)
	if len(combinations) > 0 {
		return combinations[len(combinations)-1]
	}
	return composition
}

func GetLastSoundGroup(composition []*Transformation, sound Sound) []*Transformation {
	var appendingTransformations = GetAppendingComposition(composition)
	if len(appendingTransformations) == 0 {
		return appendingTransformations
	}
	var str = Flatten(composition, NoTone|LowerCase)
	var sounds = ParseSoundsFromTonelessWord(str)
	if len(appendingTransformations) != len(sounds) {
		log.Println("The length of appending composition and sounds is not the same.")
		return appendingTransformations
	}
	var j = 0
	var i = 0
	for i = len(sounds) - 1; i >= 0; i-- {
		if sounds[i] == sound {
			j++
			if i == 0 || sounds[i-1] != sound {
				break
			}
		}
	}
	if i < 0 {
		i = 0
	}
	return appendingTransformations[i:(j + i)]
}

func GetRightMostVowels(composition []*Transformation) []*Transformation {
	return GetLastSoundGroup(composition, VowelSound)
}

func RemoveTrans(composition []*Transformation, trans *Transformation) []*Transformation {
	var transIndex = FindTransformationIndex(composition, trans)
	var t = RemoveTransIdx(composition, transIndex)
	return t
}

func RemoveTransIdx(composition []*Transformation, idx int) []*Transformation {
	if len(composition) > 0 && idx < len(composition) {
		if idx == len(composition)-1 {
			return composition[:idx]
		}
		return append(composition[:idx], composition[idx+1:]...)
	}
	return composition
}

func FindVowelByOrder(vowels []*Transformation, order int) *Transformation {
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

func FindTransformationIndex(composition []*Transformation, trans *Transformation) int {
	for i, t := range composition {
		if t == trans {
			return i
		}
	}
	return -1
}

func hasSuperWord(composition []*Transformation) bool {
	vowels := GetRightMostVowels(composition)
	if len(vowels) <= 0 {
		return false
	}
	str := Flatten(vowels, NoTone|LowerCase)
	return strings.Contains(str, "uo")
}

func isSpellingCorrect(composition []*Transformation, mode Mode) bool {
	if len(composition) <= 1 {
		return true
	}
	if mode&NoMark != 0 && mode&NoTone != 0 {
		str := Flatten(composition, NoMark|NoTone|LowerCase)
		if len([]rune(str)) <= 1 {
			return true
		}
		ok := LookupVnlDictionary(str)
		return ok
	}
	if mode&NoTone != 0 {
		str := Flatten(composition, NoTone|LowerCase)
		if len([]rune(str)) <= 1 {
			return true
		}
		ok, _ := LookupDictionary(str)
		return ok
	}
	return false
}

func isSpellingSensible(composition []*Transformation, mode Mode) bool {
	if len(composition) <= 1 {
		return true
	}
	if mode&NoTone != 0 {
		str := Flatten(composition, NoTone|LowerCase)
		if len([]rune(str)) <= 1 {
			return true
		}
		return LookupVnlDictionary(str)
	}
	return false
}

func UndoesTransformations(composition []*Transformation, applicableRules []Rule) []*Transformation {
	var result []*Transformation
	result = append(result, composition...)
	for i, trans := range result {
		for _, applicableRule := range applicableRules {
			var key = applicableRule.Key
			switch applicableRule.EffectType {
			case Appending:
				if trans.Rule.EffectType != Appending {
					continue
				}
				if key != trans.Rule.Key {
					continue
				}
				// same rule will override key and effect_on
				if trans.Rule.Effect == applicableRule.Effect {
					trans.Rule.EffectOn = AddMarkToChar(trans.Rule.EffectOn, 0)
					trans.Rule.Key = trans.Rule.EffectOn
				}
				// double typing an appending key undoes it
				if i == len(result)-1 {
					trans.IsDeleted = true
				}
				break
			case ToneTransformation:
				if trans.Rule.EffectType != ToneTransformation {
					continue
				}
				trans.IsDeleted = true
				if key == trans.Rule.Key && trans.Rule.Effect == applicableRule.Effect {
					// double typing a tone key undoes it
					// so the target will not change, the key will be appended
				} else {
					// make this tone overridable
					trans.Target = nil
				}
				break
			case MarkTransformation:
				if trans.Rule.EffectType != MarkTransformation {
					continue
				}
				if trans.Rule.EffectOn != applicableRule.EffectOn {
					continue
				}
				if key == trans.Rule.Key {
					// double typing a mark key
					trans.IsDeleted = true
				} else {
					// make this mark overridable
					trans.IsDeleted = true
					trans.Target = nil
				}
				break
			}
		}
	}
	return result
}
