/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
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

func newAppendingTrans(key rune, isUpperCase bool) *Transformation {
	return &Transformation{
		IsUpperCase: isUpperCase,
		Rule: Rule{
			Key:        key,
			EffectOn:   key,
			EffectType: Appending,
			Result:     key,
		},
	}
}

func generateAppendingTrans(rules []Rule, lowerKey rune, isUpperCase bool) *Transformation {
	for _, rule := range rules {
		if rule.Key == lowerKey && rule.EffectType == Appending {
			var _isUpperCase = isUpperCase || unicode.IsUpper(rule.EffectOn)
			rule.EffectOn = unicode.ToLower(rule.EffectOn)
			rule.Result = rule.EffectOn
			return &Transformation{
				IsUpperCase: _isUpperCase,
				Rule:        rule,
			}
		}
	}
	return newAppendingTrans(lowerKey, isUpperCase)
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

func findRootTarget(target *Transformation) *Transformation {
	if target.Target == nil {
		return target
	} else {
		return findRootTarget(target.Target)
	}
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

func getSpellingMatchResult(composition []*Transformation, mode Mode) uint8 {
	if len(composition) <= 1 {
		return FindResultMatchFull
	}
	if mode&ToneLess == 0 {
		var lastTone Tone
		// tone checking
		for i := len(composition) - 1; i >= 0; i-- {
			if composition[i].Rule.EffectType == ToneTransformation {
				lastTone = Tone(composition[i].Rule.Effect)
				break
			}
		}
		if !hasValidTone(composition, lastTone) {
			return FindResultNotMatch
		}
	}
	// spell checking
	var canvas = getCanvas(composition, VietnameseMode|LowerCase|ToneLess)
	var dictionary = mode&WithDictionary != 0
	if dictionary {
		canvas = getCanvas(composition, VietnameseMode|LowerCase)
	}
	if len(canvas) <= 1 {
		return FindResultMatchFull
	}
	return TestString(spellingTrie, canvas, dictionary)
}

func getRightMostVowels(composition []*Transformation) []*Transformation {
	return getCompositionBySound(composition, VowelSound)
}

func inComposition(composition []*Transformation, trans *Transformation) bool {
	for _, t := range composition {
		if t == trans {
			return true
		}
	}
	return false
}

func getRightMostVowelWithMarks(composition []*Transformation) []*Transformation {
	var vowels = getRightMostVowels(composition)
	var result []*Transformation
	for _, t := range composition {
		if inComposition(vowels, t) ||
			(t.Rule.EffectType == MarkTransformation && inComposition(vowels, t.Target)) {
			result = append(result, t)
		}
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
			var str = Flatten(getRightMostVowels(composition), EnglishMode|LowerCase|ToneLess|MarkLess)
			if str == "oa" || str == "oe" || str == "uy" || str == "ue" || str == "uo" {
				target = vowels[1]
			} else {
				target = vowels[0]
			}
		}
	} else if len(vowels) == 3 {
		if Flatten(vowels, EnglishMode|LowerCase|ToneLess|MarkLess) == "uye" {
			target = vowels[2]
		} else {
			target = vowels[1]
		}
	}
	return target
}

func hasValidTone(composition []*Transformation, tone Tone) bool {
	if tone == TONE_NONE || tone == TONE_ACUTE || tone == TONE_DOT {
		return true
	}
	var lastConsonants = Flatten(getCompositionBySound(composition, LastConsonantSound), EnglishMode|LowerCase)

	// These consonants have to go with ACUTE, DOT accents
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

func getLastSequence(composition []*Transformation) []*Transformation {
	for i := len(composition) - 1; i >= 0; i-- {
		if composition[i].Rule.EffectType == Appending && composition[i].Rule.Key == ' ' {
			return composition[i:]
		}
	}
	return composition
}

func getLastWord(composition []*Transformation, effectKeys []rune) []*Transformation {
	var ret []*Transformation
	for _, trans := range composition {
		ret = append(ret, trans)
		canvas := getCanvas(ret, VietnameseMode|ToneLess|LowerCase)
		if canvas == nil {
			ret = nil
			continue
		}
		var key = canvas[len(canvas)-1]
		if (IsWordBreakSymbol(key) || ('0' <= key && key <= '9')) && !inKeyList(effectKeys, key) {
			ret = nil
			continue
		}
	}
	return ret
}

func getLastSyllable(composition []*Transformation) []*Transformation {
	var ret []*Transformation
	for _, trans := range composition {
		ret = append(ret, trans)
		canvas := getCanvas(ret, VietnameseMode|ToneLess|LowerCase)
		if canvas == nil {
			ret = nil
			continue
		}
		var key = canvas[len(canvas)-1]
		if canvas == nil || IsWordBreakSymbol(key) || ('0' <= key && key <= '9') {
			ret = nil
			continue
		}
		if TestString(spellingTrie, canvas, false) == FindResultNotMatch {
			ret = []*Transformation{trans}
		}
	}
	return ret
}

func extractLastWord(composition []*Transformation, effectKeys []rune) ([]*Transformation, []*Transformation) {
	var previous, lastSyllable []*Transformation
	if len(composition) > 0 {
		var ls = getLastWord(getLastSequence(composition), effectKeys)
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

func extractLastSyllable(composition []*Transformation) ([]*Transformation, []*Transformation) {
	var previous, lastSyllable []*Transformation
	if len(composition) > 0 {
		var ls = getLastSyllable(getLastSequence(composition))
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

func findMarkTarget(composition []*Transformation, rules []Rule) (*Transformation, Rule) {
	var str = Flatten(composition, VietnameseMode)
	for i := len(composition) - 1; i >= 0; i-- {
		var trans = composition[i]
		for _, rule := range rules {
			if rule.EffectType != MarkTransformation {
				continue
			}
			if trans.Rule.Result == rule.EffectOn && rule.Effect > 0 {
				var target = findRootTarget(trans)
				if str == Flatten(append(composition, &Transformation{Target: target, Rule: rule}), VietnameseMode) {
					continue
				}
				var tmp = append(composition, &Transformation{Rule: rule, Target: target})
				if getSpellingMatchResult(tmp, ToneLess) != FindResultNotMatch {
					return target, rule
				}
			}
		}
	}
	return nil, Rule{}
}

func findTarget(composition []*Transformation, applicableRules []Rule, flags uint) (*Transformation, Rule) {
	var lastAppending = findLastAppendingTrans(composition)
	var str = Flatten(composition, VietnameseMode)
	// find tone target
	for _, applicableRule := range applicableRules {
		if applicableRule.EffectType != ToneTransformation {
			continue
		}
		var target *Transformation = nil
		if flags&EfreeToneMarking != 0 {
			if hasValidTone(composition, Tone(applicableRule.Effect)) {
				target = findToneTarget(composition, flags&EstdToneStyle != 0)
			}
		} else if lastAppending != nil && IsVowel(lastAppending.Rule.EffectOn) {
			target = lastAppending
		}
		if str == Flatten(append(composition, &Transformation{Target: target, Rule: applicableRule}), VietnameseMode) {
			continue
		}
		if Tone(applicableRule.Effect) == TONE_NONE && isFree(composition, target, ToneTransformation) &&
			FindToneFromChar(target.Rule.Result) == TONE_NONE {
			target = nil
		}
		return target, applicableRule
	}
	return findMarkTarget(composition, applicableRules)
}

func generateUndoTransformations(composition []*Transformation, rules []Rule, flags uint) []*Transformation {
	var transformations []*Transformation
	var str = Flatten(composition, VietnameseMode|ToneLess|LowerCase)
	for _, rule := range rules {
		if rule.EffectType == ToneTransformation {
			var lastAppending = findLastAppendingTrans(composition)
			var target *Transformation
			if flags&EfreeToneMarking != 0 {
				if hasValidTone(composition, Tone(rule.Effect)) {
					target = findToneTarget(composition, flags&EstdToneStyle != 0)
				}
			} else if lastAppending != nil && IsVowel(lastAppending.Rule.EffectOn) {
				target = lastAppending
			}
			if target == nil {
				continue
			}
			var trans = new(Transformation)
			trans.Target = target
			trans.Rule = Rule{
				EffectType: ToneTransformation,
				Effect:     0,
				Key:        0,
			}
			transformations = append(transformations, trans)
		} else if rule.EffectType == MarkTransformation {
			for i := len(composition) - 1; i >= 0; i-- {
				var trans = composition[i]
				if trans.Rule.Result == rule.EffectOn {
					var target = findRootTarget(trans)
					var trans = new(Transformation)
					trans.Target = target
					trans.Rule = Rule{
						Key:        0,
						EffectType: MarkTransformation,
						Effect:     0,
					}
					if str == Flatten(append(composition, trans), VietnameseMode|ToneLess|LowerCase) {
						continue
					}
					transformations = append(transformations, trans)
				}
			}
		}
	}
	return transformations
}

var regUOh_UhO_Tail = regexp.MustCompile(`\p{L}*(uơ|ưo)\p{L}+`)
var regUOh_UhO = regexp.MustCompile(`\p{L}*(uơ|ưo)\p{L}*`)
var regUhO_UhOh = regexp.MustCompile(`\p{L}*(ưo|ươ)\p{L}*`)

/**
* 1 | o + ff  ->  undo + append      -> of
* 2 | o + fs  ->  override			 -> ó
* 3 | o + fz  ->  override	         -> o
* 4 | o + z   ->  append			 -> oz
* 5 | o + f   ->  tone_grave         -> ò
* ...
**/
func generateTransformations(composition []*Transformation, applicableRules []Rule, flags uint, lowerKey rune, isUpperCase bool) []*Transformation {
	var transformations []*Transformation
	// Double typing an effect key undoes it and its effects, e.g. w + w -> w (Telex 2)
	if len(composition) > 0 {
		var rule = composition[len(composition)-1].Rule
		if rule.EffectType == Appending && rule.Key == lowerKey && rule.Key != rule.Result {
			transformations = append(transformations, &Transformation{
				Rule: Rule{
					EffectType: MarkTransformation,
					Effect:     uint8(MARK_RAW),
					Key:        0,
				},
				Target: composition[len(composition)-1],
			})
			return transformations
		}
	}
	// A target may be applied by many different transformations, e.g. o + o + w -> ơ
	if target, applicableRule := findTarget(composition, applicableRules, flags); target != nil {
		transformations = append(transformations, &Transformation{
			Rule:        applicableRule,
			Target:      target,
			IsUpperCase: isUpperCase,
		})
		var newComp = append(composition, transformations...)
		var newStr = Flatten(newComp, VietnameseMode|ToneLess|LowerCase)
		if applicableRule.EffectType == MarkTransformation && regUOh_UhO.MatchString(newStr) &&
			getSpellingMatchResult(newComp, ToneLess) == FindResultMatchPrefix {
			// Implement the uow typing shortcut by creating a virtual
			// Mark_HORN rule that targets 'u' or 'o'.
			if target, virtualRule := findTarget(newComp, applicableRules, flags); target != nil {
				virtualRule.Key = 0
				transformations = append(transformations, &Transformation{
					Rule:   virtualRule,
					Target: target,
				})
			}
		}
	} else {
		// Implement ươ/ưo + o -> uô
		if regUhO_UhOh.MatchString(Flatten(composition, VietnameseMode|ToneLess|LowerCase)) {
			var vowels = getRightMostVowelWithMarks(composition)
			var trans = &Transformation{
				Target: vowels[0],
				Rule: Rule{
					EffectType: MarkTransformation,
					Key:        0,
					Effect:     uint8(MARK_NONE),
				},
			}
			if target, applicableRule := findTarget(append(composition, trans), applicableRules, flags); target != nil && target != vowels[0] {
				transformations = append(transformations, trans)
				transformations = append(transformations, &Transformation{
					Rule:        applicableRule,
					Target:      target,
					IsUpperCase: isUpperCase,
				})
				return transformations
			}
		}
		if undoTrans := generateUndoTransformations(composition, applicableRules, flags); len(undoTrans) > 0 {
			// If an effect key can't find its target, it tries to undo its brothers, e.g. ươ + w -> uow
			transformations = append(transformations, undoTrans...)
			transformations = append(transformations, newAppendingTrans(lowerKey, isUpperCase))
		}
	}
	return transformations
}

func generateFallbackTransformations(composition []*Transformation, applicableRules []Rule, lowerKey rune, isUpperCase bool) []*Transformation {
	var transformations []*Transformation
	var trans = generateAppendingTrans(applicableRules, lowerKey, isUpperCase)
	transformations = append(transformations, trans)
	for _, appendedRule := range trans.Rule.AppendedRules {
		var _isUpperCase = isUpperCase || unicode.IsUpper(appendedRule.EffectOn)
		appendedRule.Key = 0 // this is a virtual key
		appendedRule.EffectOn = unicode.ToLower(appendedRule.EffectOn)
		appendedRule.Result = appendedRule.EffectOn
		transformations = append(transformations, &Transformation{
			Rule:        appendedRule,
			IsUpperCase: _isUpperCase,
		})
	}
	if composition != nil && trans.Rule.Key != trans.Rule.Result && getSpellingMatchResult(append(composition, transformations...), ToneLess) == FindResultNotMatch {
		transformations = []*Transformation{newAppendingTrans(lowerKey, isUpperCase)}
	}
	return transformations
}

func breakComposition(composition []*Transformation) []*Transformation {
	var result []*Transformation
	for _, trans := range composition {
		if trans.Rule.Key == 0 {
			continue
		}
		result = append(result, newAppendingTrans(trans.Rule.Key, trans.IsUpperCase))
	}
	return result
}

func refreshLastToneTarget(composition []*Transformation, stdStyle bool) []*Transformation {
	var transformations []*Transformation
	var rightmostVowels = getRightMostVowels(composition)
	var lastToneTrans = getLastToneTransformation(composition)
	if len(rightmostVowels) == 0 || lastToneTrans == nil {
		return nil
	}
	var newToneTarget = findToneTarget(composition, stdStyle)
	if lastToneTrans.Target != newToneTarget {
		lastToneTrans.Target = newToneTarget
		transformations = append(transformations, &Transformation{
			Target: lastToneTrans.Target,
			Rule: Rule{
				Key:        0,
				EffectType: ToneTransformation,
				Effect:     uint8(TONE_NONE),
			},
		})
		var overrideRule = lastToneTrans.Rule
		overrideRule.Key = 0
		transformations = append(transformations, &Transformation{
			Target: newToneTarget,
			Rule:   overrideRule,
		})
	}
	return transformations
}
