/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LISENCE>.
 */
package bamboo

import (
	"regexp"
	"strings"
	"unicode"
)

func findAppendingRule(rules []Rule, key rune) Rule {
	var result Rule
	result.EffectType = Appending
	result.Key = key
	result.EffectOn = key
	var applicableRules []Rule
	for _, inputRule := range rules {
		if inputRule.Key == key {
			applicableRules = append(applicableRules, inputRule)
		}
	}
	for _, applicableRule := range applicableRules {
		if applicableRule.EffectType == Appending {
			result.EffectOn = applicableRule.EffectOn
			result.AppendedRules = applicableRule.AppendedRules
			return result
		}
	}
	return result
}

func findLastAppendingTrans(composition []*Transformation) *Transformation {
	for i := len(composition) - 1; i >= 0; i-- {
		var trans = composition[i]
		if trans.Rule.EffectType == Appending {
			return trans
		}
	}
	return nil
}

func findNextAppendingTransformation(composition []*Transformation, trans *Transformation) (*Transformation, bool) {
	fromIndex := findTransformationIndex(composition, trans)
	if fromIndex == -1 {
		return nil, false
	}
	var nextAppendingTrans *Transformation
	found := false
	for i := fromIndex + 1; int(i) < len(composition); i++ {
		if composition[i].Rule.EffectType == Appending {
			nextAppendingTrans = composition[i]
			found = true
		}
	}
	return nextAppendingTrans, found
}

func createAppendingTrans(key rune, isUpperCase bool) *Transformation {
	return &Transformation{
		IsUpperCase: isUpperCase,
		Rule: Rule{
			Key:        unicode.ToLower(key),
			EffectOn:   unicode.ToLower(key),
			EffectType: Appending,
		},
	}
}
func getAppendingComposition(composition []*Transformation) []*Transformation {
	var appendingTransformations []*Transformation
	for _, trans := range composition {
		if trans.Rule.EffectType == Appending {
			appendingTransformations = append(appendingTransformations, trans)
		}
	}
	return appendingTransformations
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

func findToneTarget(composition []*Transformation, stdStyle bool) *Transformation {
	if len(composition) == 0 {
		return nil
	}
	var target *Transformation
	var vowels = getRightMostVowels(composition)
	if len(vowels) == 1 {
		target = vowels[0]
	} else if len(vowels) == 2 && stdStyle {
		var str = Flatten(getRightMostVowelWithMarks(composition), NoTone|LowerCase)
		var chars = []rune(str)
		if ohPos := findIndexRune(chars, 'ơ'); ohPos > 0 {
			target = vowels[ohPos]
		} else if ehPos := findIndexRune(chars, 'ê'); ehPos > 0 {
			target = vowels[ehPos]
		} else {
			if _, found := findNextAppendingTransformation(composition, vowels[1]); found {
				target = vowels[1]
			} else {
				target = vowels[0]
			}
		}
	} else if len(vowels) == 2 {
		var str = Flatten(getRightMostVowels(composition), EnglishMode|LowerCase)
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
	var lastConsonants = Flatten(getCompositionBySound(composition, LastConsonantSound), EnglishMode|LowerCase)

	// These consonants can only go with ACUTE, DOT accents
	var dotWithConsonants = []string{"c", "k", "p", "t", "ch"}
	for _, s := range dotWithConsonants {
		if s == lastConsonants {
			return false
		}
	}
	return true
}

func getLastToneTransformation(composition []*Transformation) *Transformation {
	for i := len(composition) - 1; i >= 0; i-- {
		var t = composition[i]
		if t.Rule.EffectType == ToneTransformation && t.Target != nil {
			return t
		}
	}
	return nil
}

func shouldRefreshLastToneTarget(transformations []*Transformation) bool {
	var vowels = getRightMostVowels(transformations)
	if len(vowels) <= 0 {
		return false
	}
	return len(transformations) > 0 && transformations[len(transformations)-1].Rule.EffectType == Appending
}

func refreshLastToneTarget(transformations []*Transformation) []*Transformation {
	var composition []*Transformation
	composition = append(composition, transformations...)
	var rightmostVowels = getRightMostVowels(composition)
	var lastToneTrans = getLastToneTransformation(composition)
	if len(rightmostVowels) == 0 || lastToneTrans == nil {
		return composition
	}
	var newToneTarget = findToneTarget(composition, true)
	if lastToneTrans.Target != newToneTarget {
		lastToneTrans.Target = newToneTarget
	}
	return composition
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

/***** BEGIN SIDE-EFFECT METHODS ******/

func removeTrans(composition []*Transformation, trans *Transformation) []*Transformation {
	var transIndex = findTransformationIndex(composition, trans)
	if transIndex < 0 {
		return composition
	}
	return removeTransIdx(composition, transIndex)
}

func removeTransIdx(composition []*Transformation, idx int) []*Transformation {
	if len(composition) > 0 && idx >= 0 && idx < len(composition) {
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
