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
	"regexp"
	"strings"
)

var regUho = regexp.MustCompile(`uơ\p{L}*`)
var regUoh = regexp.MustCompile(`ưo\p{L}*`)
var regUoSuffix = regexp.MustCompile(`^(h|th|kh)uo$`)

func findMissingRuleForUo(composition []*Transformation, isSuperKey bool) (Rule, bool) {
	var rule Rule
	if len(composition) < 2 {
		return rule, false
	}
	var target rune
	var full = strings.ToLower(Flatten(composition, NoTone|LowerCase))

	if !isSuperKey {
		if regUho.MatchString(full) {
			target = 'u'
		}
		if regUoh.MatchString(full) {
			target = 'o'
		}
	} else {
		if regUoSuffix.MatchString(full) {
			return rule, false
		}
		var vowels = getRightMostVowelWithMarks(composition)
		var str = Flatten(vowels, NoTone|LowerCase)
		if strings.Contains(str, "uo") {
			target = 'o'
		}
	}
	if target > 0 {
		rule = Rule{
			Key:        rune(0),
			EffectType: MarkTransformation,
			Effect:     uint8(MARK_HORN),
			EffectOn:   target,
		}
		return rule, true
	}
	return rule, false
}

func findIndexRune(chars []rune, r rune) int {
	for i, c := range chars {
		if c == r {
			return i
		}
	}
	return -1
}

func isFree(composition []*Transformation, trans *Transformation, effectType EffectType) bool {
	for _, t := range composition {
		if t.Target == trans && t.Rule.EffectType == effectType {
			return false
		}
	}
	return true
}

func findTransformationIndex(composition []*Transformation, trans *Transformation) int {
	for i, t := range composition {
		if t == trans {
			return i
		}
	}
	return -1
}

func hasSuperWord(composition []*Transformation) bool {
	vowels := getRightMostVowels(composition)
	if len(vowels) <= 0 {
		return false
	}
	str := Flatten(vowels, NoTone|LowerCase)
	return strings.Contains(str, "uo")
}

/***** BEGIN SIDE-EFFECT METHODS ******/

func removeTrans(composition []*Transformation, trans *Transformation) []*Transformation {
	var transIndex = findTransformationIndex(composition, trans)
	var t = removeTransIdx(composition, transIndex)
	return t
}

func removeTransIdx(composition []*Transformation, idx int) []*Transformation {
	if len(composition) > 0 && idx < len(composition) {
		if idx == len(composition)-1 {
			return composition[:idx]
		}
		return append(composition[:idx], composition[idx+1:]...)
	}
	return composition
}

func undoesTransformations(composition []*Transformation, applicableRules []Rule) []*Transformation {
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

func freeComposition(composition []*Transformation) []*Transformation {
	var result []*Transformation
	result = append(result, composition...)
	for i, trans := range composition {
		if trans.IsDeleted {
			result = removeTransIdx(result, i)
		}
	}
	return result
}

/***** END SIDE-EFFECT METHODS ******/
