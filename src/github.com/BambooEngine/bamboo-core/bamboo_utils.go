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
	"log"
	"regexp"
	"unicode"
)

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
			Key:        key,
			EffectOn:   key,
			EffectType: Appending,
		},
	}
}

func createAppendingComposition(rules []Rule, lowerKey rune, isUpperCase bool) []*Transformation {
	var transformations []*Transformation
	for _, rule := range rules {
		if rule.Key == lowerKey && rule.EffectType == Appending {
			var _isUpperCase = isUpperCase || unicode.IsUpper(rule.EffectOn)
			rule.EffectOn = unicode.ToLower(rule.EffectOn)
			transformations = append(transformations, &Transformation{
				IsUpperCase: _isUpperCase,
				Rule:        rule,
			})
			for _, appendedRule := range rule.AppendedRules {
				var _isUpperCase = isUpperCase || unicode.IsUpper(appendedRule.EffectOn)
				appendedRule.Key = 0 // this is a virtual key
				appendedRule.EffectOn = unicode.ToLower(appendedRule.EffectOn)
				transformations = append(transformations, &Transformation{
					Rule:        appendedRule,
					IsUpperCase: _isUpperCase,
				})
			}
		}
	}
	if len(transformations) == 0 {
		transformations = append(transformations, createAppendingTrans(lowerKey, isUpperCase))
	}
	return transformations
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
	// a target is not applied mark transformation twice
	if !isFree(composition, trans.Target, MarkTransformation) {
		return false
	}
	var soundMap = GetSoundMap(composition)
	targetSound, found := soundMap[trans.Target]
	if !found {
		return false
	}
	// the sound of target does not match the sound of the trans that effect on
	if IsVowel(trans.Rule.EffectOn) && targetSound != VowelSound {
		return false
	}
	if targetSound == VowelSound {
		var vowels = getRightMostVowelWithMarks(append(composition, trans))
		if getSpellingMatchResult(vowels, ToneLess, false) == FindResultNotMatch {
			return false
		}
	}
	return true
}

// only appending trans has sound
func getCombinationWithSound(composition []*Transformation) ([]*Transformation, []Sound) {
	var lastComb = getAppendingComposition(composition)
	if len(lastComb) <= 0 {
		return lastComb, nil
	}
	var str = Flatten(lastComb, VietnameseMode|ToneLess|LowerCase)
	if TestString(spellingTrie, []rune(str), false) != FindResultNotMatch {
		return lastComb, ParseSoundsFromWord(str)
	}
	return lastComb, ParseSoundsFromWord(str)
}

func getCompositionBySound(composition []*Transformation, sound Sound) []*Transformation {
	var lastComb, sounds = getCombinationWithSound(composition)
	if len(lastComb) != len(sounds) {
		log.Println("Something is wrong with the length of sounds")
		return lastComb
	}
	var ret []*Transformation
	for i, s := range sounds {
		if s == sound {
			ret = append(ret, lastComb[i])
		}
	}
	return ret
}

func getSpellingMatchResult(composition []*Transformation, mode Mode, deepSearch bool) uint8 {
	if len(composition) <= 0 {
		return FindResultMatchFull
	}
	if mode&ToneLess != 0 {
		str := Flatten(composition, ToneLess|LowerCase)
		var chars = []rune(str)
		if len(chars) <= 1 {
			return FindResultMatchFull
		}
		return TestString(spellingTrie, chars, deepSearch)
	}
	return FindResultNotMatch
}

func isSpellingCorrect(composition []*Transformation, mode Mode) bool {
	res := getSpellingMatchResult(composition, mode, false)
	return res == FindResultMatchFull
}

func GetSoundMap(composition []*Transformation) map[*Transformation]Sound {
	var soundMap = map[*Transformation]Sound{}
	var lastComb, sounds = getCombinationWithSound(composition)
	if len(sounds) <= 0 || len(sounds) != len(lastComb) {
		log.Println("Something is wrong with the length of sounds")
		return soundMap
	}
	for i, trans := range lastComb {
		soundMap[trans] = sounds[i]
	}
	return soundMap
}

func getRightMostVowels(composition []*Transformation) []*Transformation {
	return getCompositionBySound(composition, VowelSound)
}

func getRightMostVowelWithMarks(composition []*Transformation) []*Transformation {
	var vowels = getRightMostVowels(composition)
	return addMarksToComposition(composition, vowels)
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
		var str = Flatten(getRightMostVowelWithMarks(composition), ToneLess|LowerCase)
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
		if _, found := findNextAppendingTransformation(composition, vowels[1]); found {
			target = vowels[1]
		} else {
			var str = Flatten(getRightMostVowels(composition), EnglishMode|LowerCase)
			if str == "oa" || str == "oe" || str == "uy" || str == "ue" || str == "uo" {
				target = vowels[1]
			} else {
				target = vowels[0]
			}
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

func isTransformationForUoMissed(composition []*Transformation) bool {
	return len(composition) > 0 &&
		hasSuperWord(composition) &&
		getSpellingMatchResult(composition, ToneLess, false) == FindResultMatchPrefix
}

func refreshLastToneTarget(transformations []*Transformation, stdStyle bool) []*Transformation {
	var composition []*Transformation
	composition = append(composition, transformations...)
	var rightmostVowels = getRightMostVowels(composition)
	var lastToneTrans = getLastToneTransformation(composition)
	if len(rightmostVowels) == 0 || lastToneTrans == nil {
		return composition
	}
	var newToneTarget = findToneTarget(composition, stdStyle)
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

var regUhOh = regexp.MustCompile(`\p{L}*(uơ|ưo)\p{L}*`)

func hasSuperWord(composition []*Transformation) bool {
	str := Flatten(composition, ToneLess|LowerCase)
	return regUhOh.MatchString(str)
}

func extractLastSyllable(composition []*Transformation) ([]*Transformation, []*Transformation) {
	var previous, lastSyllable []*Transformation
	if len(composition) > 0 {
		var ls = getLastSyllable(getLastWord(composition, nil))
		if len(ls) > 0 {
			var idx = findTransformationIndex(composition, ls[0])
			if idx > 0 {
				previous = composition[:idx]
			}
			lastSyllable = ls
		} else {
			previous = composition
		}
	}
	return lastSyllable, previous
}

func findTargetFromKey(composition []*Transformation, applicableRules []Rule, flags uint) (*Transformation, Rule) {
	var lastAppending = findLastAppendingTrans(composition)
	for _, applicableRule := range applicableRules {
		var target *Transformation = nil
		if applicableRule.EffectType == MarkTransformation {
			return findMarkTarget(composition, applicableRules)
		} else if applicableRule.EffectType == ToneTransformation {
			if flags&EfreeToneMarking != 0 {
				if hasValidTone(composition, Tone(applicableRule.Effect)) {
					target = findToneTarget(composition, flags&EstdToneStyle != 0)
					if !isFree(composition, target, ToneTransformation) {
						target = nil
					}
				}
			} else if lastAppending != nil && IsVowel(lastAppending.Rule.EffectOn) {
				target = lastAppending
			}
		}
		if target != nil {
			return target, applicableRule
		}
	}
	return nil, Rule{}
}

// If none of the applicable_rules can actually be applied then this new
// transformation fall-backs to an APPENDING one.
func generateTransformations(composition []*Transformation, applicableRules []Rule, flags uint, key rune, isUpperCase bool) []*Transformation {
	var transformations []*Transformation
	if target, applicableRule := findTargetFromKey(composition, applicableRules, flags); target != nil {
		transformations = append(transformations, &Transformation{
			Rule:        applicableRule,
			Target:      target,
			IsUpperCase: isUpperCase,
		})
	} else {
		transformations = append(transformations, createAppendingComposition(applicableRules, key, isUpperCase)...)
	}
	return transformations
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
	for i, trans := range composition {
		if trans.IsDeleted {
			composition = removeTransIdx(composition, i)
		}
	}
	return composition
}

/***** END SIDE-EFFECT METHODS ******/
