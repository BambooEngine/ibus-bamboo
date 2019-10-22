/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LISENCE>.
 */
package bamboo

import (
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

func filterAppendingComposition(composition []*Transformation) []*Transformation {
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

func isValid(composition []*Transformation, inputIsFullComplete bool) bool {
	if len(composition) <= 1 {
		return true
	}
	// last tone checking
	for i := len(composition) - 1; i >= 0; i-- {
		if composition[i].Rule.EffectType == ToneTransformation {
			var lastTone = Tone(composition[i].Rule.Effect)
			if !hasValidTone(composition, lastTone) {
				return false
			}
			break
		}
	}
	// spell checking
	var fc, vo, lc = extractCvcTrans(composition)
	var flattenMode = VietnameseMode | LowerCase | ToneLess
	return isValidCVC(Flatten(fc, flattenMode), Flatten(vo, flattenMode), Flatten(lc, flattenMode), inputIsFullComplete)
}

func getRightMostVowels(composition []*Transformation) []*Transformation {
	var _, vo, _ = extractCvcTrans(composition)
	return vo
}

func findToneTarget(composition []*Transformation, stdStyle bool) *Transformation {
	if len(composition) == 0 {
		return nil
	}
	var target *Transformation
	var _, vo, lc = extractCvcTrans(composition)
	var vowels = filterAppendingComposition(vo)
	if len(vowels) == 1 {
		target = vowels[0]
	} else if len(vowels) == 2 && stdStyle {
		for _, trans := range vo {
			if trans.Rule.Result == 'ơ' || trans.Rule.Result == 'ê' {
				if trans.Target != nil {
					target = trans.Target
				} else {
					target = trans
				}
			}
		}
		if target == nil {
			if len(lc) > 0 {
				target = vowels[1]
			} else {
				target = vowels[0]
			}
		}
	} else if len(vowels) == 2 {
		if len(lc) > 0 {
			target = vowels[1]
		} else {
			var str = Flatten(vowels, EnglishMode|LowerCase|ToneLess|MarkLess)
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
	var _, _, lc = extractCvcTrans(composition)
	if len(lc) == 0 {
		return true
	}
	var lastConsonants = Flatten(lc, EnglishMode|LowerCase)

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
			if i < len(composition)-1 {
				return composition[i+1:]
			} else {
				return nil
			}
		}
	}
	return composition
}

func getLastWord(composition []*Transformation, effectKeys []rune) []*Transformation {
	var ret []*Transformation
	for _, trans := range composition {
		ret = append(ret, trans)
		var canvas = getCanvas(ret, VietnameseMode)
		if len(canvas) == 0 {
			continue
		}
		var key = canvas[len(canvas)-1]
		if IsWordBreakSymbol(key) && !inKeyList(effectKeys, key) {
			ret = nil
			continue
		}
	}
	return ret
}

func getLastSyllable(composition []*Transformation) []*Transformation {
	var ret []*Transformation
	if isValid(composition, false) {
		return composition
	}
	for _, trans := range composition {
		ret = append(ret, trans)
		var canvas = getCanvas(ret, VietnameseMode)
		if len(canvas) == 0 {
			continue
		}
		var key = canvas[len(canvas)-1]
		if IsWordBreakSymbol(key) {
			ret = nil
			continue
		}
		if !isValid(ret, false) {
			ret = []*Transformation{trans}
		}
	}
	return ret
}

func extractAtomicTrans(composition, last []*Transformation, lastIsVowel bool) ([]*Transformation, []*Transformation) {
	if len(composition) == 0 {
		return composition, last
	}
	var tmp = composition[len(composition)-1]
	if tmp != nil && tmp.Target == nil && lastIsVowel != IsVowel(tmp.Rule.Result) {
		return composition, last
	}
	return extractAtomicTrans(composition[:len(composition)-1], append([]*Transformation{composition[len(composition)-1]}, last...), lastIsVowel)
}

/*
	Separate a string into smaller parts: first consonant (or head), vowel,
	last consonant (if any).
*/
func extractCvcAppendingTrans(composition []*Transformation) ([]*Transformation, []*Transformation, []*Transformation) {
	head, lastConsonant := extractAtomicTrans(composition, nil, false)
	firstConsonant, vowel := extractAtomicTrans(head, nil, true)
	if len(lastConsonant) > 0 && len(vowel) == 0 && len(firstConsonant) == 0 {
		firstConsonant = lastConsonant
		vowel = nil
		lastConsonant = nil
	}

	// 'gi' and 'qu' are considered qualified consonants.
	// We want something like this:
	//     ['g', 'ia', ''] -> ['gi', 'a', '']
	//     ['q', 'ua', ''] -> ['qu', 'a', '']
	// except:
	//     ['g', 'ie', 'ng'] -> ['g', 'ie', 'ng']
	if len(firstConsonant) == 1 && len(vowel) > 0 && ((firstConsonant[0].Rule.Result == 'g' && vowel[0].Rule.Result == 'i' && len(vowel) > 1 &&
		!(vowel[1].Rule.Result == 'e' && len(lastConsonant) > 0)) ||
		(firstConsonant[0].Rule.Result == 'q' && vowel[0].Rule.Result == 'u')) {
		firstConsonant = append(firstConsonant, vowel[0])
		vowel = vowel[1:]
	}
	return firstConsonant, vowel, lastConsonant
}

func extractCvcTrans(composition []*Transformation) ([]*Transformation, []*Transformation, []*Transformation) {
	var transMap = map[*Transformation][]*Transformation{}
	var appendingList []*Transformation
	for _, trans := range composition {
		if trans.Target == nil {
			appendingList = append(appendingList, trans)
		} else {
			transMap[trans.Target] = append(transMap[trans.Target], trans)
		}
	}
	var fc, vo, lc = extractCvcAppendingTrans(appendingList)
	for _, t := range fc {
		fc = append(fc, transMap[t]...)
	}
	for _, t := range vo {
		vo = append(vo, transMap[t]...)
	}
	for _, t := range lc {
		lc = append(lc, transMap[t]...)
	}
	return fc, vo, lc
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
				if isValid(tmp, false) {
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

var regUOh_UhO_Tail = regexp.MustCompile(`(uơ|ưo)\p{L}+`)
var regUOh_UhO = regexp.MustCompile(`(\p{L}*)(uơ|ưo)(\p{L}*)`)
var regUhO_UhOh = regexp.MustCompile(`(ưo|ươ)`)

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
		if applicableRule.EffectType == MarkTransformation && regUOh_UhO.MatchString(newStr) {
			var subs = regUOh_UhO.FindStringSubmatch(newStr)
			var lcIndexes = lookup(firstConsonantSeqs, subs[1], false, true)
			if (len(lcIndexes) == 1 && lcIndexes[0] != 1) || subs[3] != "" {
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
		}
	} else {
		// Implement ươ/ưo + o -> uô
		if regUhO_UhOh.MatchString(Flatten(composition, VietnameseMode|ToneLess|LowerCase)) {
			var vowels = filterAppendingComposition(getRightMostVowels(composition))
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
			// If an effect key can't find its target, it tries to undo its effects, e.g. ươ + w -> uow
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
	if composition != nil && IsPunctuationMark(trans.Rule.Key) && !isValid(append(composition, transformations...), false) {
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
	if rightmostVowels == nil || lastToneTrans == nil {
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
