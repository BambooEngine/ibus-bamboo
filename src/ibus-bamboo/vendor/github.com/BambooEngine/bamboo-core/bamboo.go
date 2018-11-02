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
	"unicode"
)

type Mode uint

const (
	VietnameseMode Mode = 1 << iota
	EnglishMode
	NoTone
	NoMark
	LowerCase
)

type Sound uint

const (
	NoSound             Sound = iota << 0
	FirstConsonantSound Sound = iota
	VowelSound          Sound = iota
	LastConsonantSound  Sound = iota
)

const (
	EfreeToneMarking uint = 1 << iota
	EstdToneStyle
	EspellCheckEnabled
	EmarcoEnabled
	EautoNonVnRestore
	EautoCorrect
	EfastCommitting
	EstdFlags = EfreeToneMarking | EstdToneStyle | EautoCorrect | EspellCheckEnabled
)

type Transformation struct {
	Rule        Rule
	Target      *Transformation // For Tone/Mark transformation
	IsUpperCase bool
	Dest        uint // For Appending, a pointer to the char in the flattened string made by this Trans
}

type IEngine interface {
	SetFlag(uint)
	GetInputMethod() InputMethod
	ProcessChar(rune)
	ProcessString(string)
	GetProcessedString(Mode) string
	IsSpellingCorrect(Mode) bool
	IsSpellingSensible(Mode) bool
	Reset()
	RemoveLastChar()
}

type BambooEngine struct {
	composition []*Transformation
	inputMethod InputMethod
	flags       uint
}

func NewEngine(im string, flag uint) IEngine {
	inputMethod, found := InputMethods[im]
	if !found {
		panic("The input method is not supported")
	}
	engine := BambooEngine{
		inputMethod: inputMethod,
		flags:       flag,
	}
	return &engine
}

func (e *BambooEngine) GetInputMethod() InputMethod {
	return e.inputMethod
}

func (e *BambooEngine) SetFlag(flag uint) {
	e.flags = flag
}

func (e *BambooEngine) GetFlag(flag uint) uint {
	return e.flags
}

func (e *BambooEngine) isSuperKey(chr rune) bool {
	for _, key := range e.inputMethod.SuperKeys {
		if key == chr {
			return true
		}
	}
	return false
}

func (e *BambooEngine) isToneKey(chr rune) bool {
	for _, key := range e.inputMethod.ToneKeys {
		if key == chr {
			return true
		}
	}
	return false
}

func (e *BambooEngine) isEffectiveKey(chr rune) bool {
	for _, key := range e.inputMethod.Keys {
		if key == chr {
			return true
		}
	}
	return false
}

func (e *BambooEngine) isFree(trans *Transformation, effectType EffectType) bool {
	for _, t := range e.composition {
		if t.Target == trans && t.Rule.EffectType == effectType {
			return false
		}
	}
	return true
}

func (e *BambooEngine) isCharFree(c rune, effectType EffectType) bool {
	for _, t := range e.composition {
		if t.Rule.EffectOn == c && t.Rule.EffectType == effectType {
			return false
		}
	}
	return true
}

func (e *BambooEngine) findTargetForKey(chr rune) (*Transformation, Rule) {
	var applicableRules []Rule
	for _, inputRule := range e.inputMethod.Rules {
		if inputRule.Key == chr {
			applicableRules = append(applicableRules, inputRule)
		}
	}
	var lastAppending = FindLastAppendingTrans(e.composition)
	for _, applicableRule := range applicableRules {
		var target *Transformation = nil
		if applicableRule.EffectType == MarkTransformation {
			return FindMarkTarget(e.composition, applicableRules)
		} else if applicableRule.EffectType == ToneTransformation {
			if e.flags&EfreeToneMarking != 0 {
				if hasValidTone(e.composition, Tone(applicableRule.Effect)) {
					target = FindToneTarget(e.composition, e.flags&EstdToneStyle != 0)
				}
			} else if lastAppending != nil && isVowel(lastAppending.Rule.EffectOn) {
				target = lastAppending
			}
		}
		if target != nil {
			return target, applicableRule
		}
	}
	return nil, Rule{}
}

func (e *BambooEngine) createCompositionForRule(rule Rule, isUpperKey bool) []*Transformation {
	var transformations []*Transformation
	var trans = new(Transformation)
	trans.Rule = rule
	trans.IsUpperCase = isUpperKey
	if target, applicableRule := e.findTargetForKey(rule.Key); target != nil {
		trans.Rule = applicableRule
		trans.Target = target
	}
	transformations = append(transformations, trans)
	for _, appendedRule := range trans.Rule.AppendedRules {
		transformations = append(transformations, &Transformation{Rule: appendedRule})
	}
	return transformations
}

func (e *BambooEngine) IsSpellingCorrect(mode Mode) bool {
	return isSpellingCorrect(e.composition, mode)
}

func (e *BambooEngine) IsSpellingSensible(mode Mode) bool {
	return isSpellingSensible(e.composition, mode)
}

func (e *BambooEngine) createCompositionForKey(chr rune) []*Transformation {
	var isUpperCase bool
	if unicode.IsUpper(chr) {
		isUpperCase = true
	}
	chr = unicode.ToLower(chr)
	var transformations []*Transformation
	transformations = e.createCompositionForRule(FindAppendingRule(e.inputMethod.Rules, chr), isUpperCase)
	return transformations
}

func (e *BambooEngine) GetRawString() string {
	var seq []rune
	for _, t := range e.composition {
		seq = append(seq, t.Rule.Key)
	}
	return string(seq)
}

func (e *BambooEngine) GetProcessedString(mode Mode) string {
	return Flatten(e.composition, mode)
}

/***** BEGIN SIDE-EFFECT METHODS ******/

func (e *BambooEngine) refreshLastToneTarget() {
	// Refresh the tone position of the rightmost vowels
	var rightmostVowels = GetRightMostVowels(e.composition)
	if len(rightmostVowels) <= 0 {
		return
	}
	var rightmostVowelPos = FindTransPos(e.composition, rightmostVowels[0])
	for i := len(e.composition) - 1; i >= 0; i-- {
		trans := e.composition[i]
		if trans.Rule.EffectType == ToneTransformation {
			var newTarget = FindToneTarget(e.composition, e.flags&EstdToneStyle != 0)
			var tonePos = FindTransPos(e.composition, trans)
			if tonePos > rightmostVowelPos {
				trans.Target = newTarget
				break
			}
		}
	}
}

func (e *BambooEngine) ProcessChar(key rune) {
	if e.flags&EautoCorrect != 0 && (e.isSuperKey(key) || (!e.isToneKey(key) && hasSuperWord(e.composition))) {
		if missingRule, found := FindMissingRuleForUo(e.composition, e.isSuperKey(key)); found {
			var targets = FindMarkTargets(e.composition, missingRule)
			if len(targets) > 0 {
				virtualTrans := &Transformation{
					Rule:   missingRule,
					Target: targets[len(targets)-1],
				}
				e.composition = append(e.composition, virtualTrans)
			} else {
				log.Println("Cannot find targets for the missing rule for uo")
			}
		}
	}
	transformations := e.createCompositionForKey(key)
	e.composition = append(e.composition, transformations...)

	// double-keystroke to revert last effect, e.g in telex ww->w, ff->f, ss->s,...
	if e.isEffectiveKey(key) && haveDoubledKeystroke(e.composition, key) {
		e.composition = e.composition[0 : len(e.composition)-2]
		e.composition = append(e.composition, &Transformation{
			Rule: Rule{Key: key, EffectType: Appending, EffectOn: key},
		})
	}

	/**
	* Sometimes, a tone's position in a previous state must be changed to fit the new state
	*
	* e.g.
	* prev state: chuyenr -> chuỷen
	* this state: chuyenre -> chuyển
	**/
	if e.flags&EstdToneStyle != 0 && shouldRefreshLastToneTarget(e.composition) {
		e.refreshLastToneTarget()
	}
}

func (e *BambooEngine) ProcessString(str string) {
	for _, chr := range []rune(str) {
		e.ProcessChar(chr)
	}
}

func (e *BambooEngine) Reset() {
	e.composition = []*Transformation{}
}

func (e *BambooEngine) RemoveLastChar() {
	var lastAppending = FindLastAppendingTrans(e.composition)
	var transformations = GetTransformationsTargetTo(e.composition, lastAppending)
	for _, trans := range append(transformations, lastAppending) {
		e.composition = RemoveTrans(e.composition, trans)
	}
}

/***** END SIDE-EFFECT METHODS ******/
