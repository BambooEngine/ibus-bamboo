/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENSE>.
 */

// Package bamboo implements text processing for Vietnamese
package bamboo

import (
	"unicode"
)

type Mode uint

const (
	VietnameseMode Mode = 1 << iota
	EnglishMode
	ToneLess
	MarkLess
	LowerCase
	FullText
	PunctuationMode
	InReverseOrder
)

const (
	EfreeToneMarking uint = 1 << iota
	EstdToneStyle
	EautoCorrectEnabled
	EstdFlags = EfreeToneMarking | EstdToneStyle | EautoCorrectEnabled
)

type Transformation struct {
	Rule        Rule
	Target      *Transformation
	IsUpperCase bool
}

type IEngine interface {
	SetFlag(uint)
	GetInputMethod() InputMethod
	ProcessKey(rune, Mode)
	ProcessString(string, Mode)
	GetProcessedString(Mode) string
	IsValid(bool) bool
	CanProcessKey(rune) bool
	RemoveLastChar(bool)
	RestoreLastWord(bool)
	Reset()
}

type BambooEngine struct {
	composition []*Transformation
	inputMethod InputMethod
	flags       uint
}

func NewEngine(inputMethod InputMethod, flag uint) IEngine {
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

func (e *BambooEngine) IsValid(inputIsFullComplete bool) bool {
	var _, last = extractLastWord(e.composition, e.GetInputMethod().Keys)
	return isValid(last, inputIsFullComplete)
}

func (e *BambooEngine) GetProcessedString(mode Mode) string {
	var tmp []*Transformation
	if mode&FullText != 0 {
		tmp = e.composition
	} else if mode&PunctuationMode != 0 {
		_, tmp = extractLastWordWithPunctuationMarks(e.composition, e.inputMethod.Keys)
		return Flatten(tmp, VietnameseMode)
	} else {
		_, tmp = extractLastWord(e.composition, e.inputMethod.Keys)
	}
	return Flatten(tmp, mode)
}

func (e *BambooEngine) getApplicableRules(key rune) []Rule {
	var applicableRules []Rule
	for _, inputRule := range e.inputMethod.Rules {
		if inputRule.Key == unicode.ToLower(key) {
			applicableRules = append(applicableRules, inputRule)
		}
	}
	return applicableRules
}

func (e *BambooEngine) findTargetByKey(composition []*Transformation, key rune) (*Transformation, Rule) {
	return findTarget(composition, e.getApplicableRules(key), e.flags)
}

func (e *BambooEngine) CanProcessKey(key rune) bool {
	return canProcessKey(key, e.inputMethod.Keys)
}

func (e *BambooEngine) generateTransformations(composition []*Transformation, lowerKey rune, isUpperCase bool) []*Transformation {
	var transformations = generateTransformations(composition, e.getApplicableRules(lowerKey), e.flags, lowerKey, isUpperCase)
	if transformations == nil {
		// If none of the applicable_rules can actually be applied then this new
		// transformation fall-backs to an APPENDING one.
		transformations = generateFallbackTransformations(composition, e.getApplicableRules(lowerKey), lowerKey, isUpperCase)
		var newComposition = append(composition, transformations...)

		// Implement the uwo+ typing shortcut by creating a virtual
		// Mark.HORN rule that targets 'u' or 'o'.
		if virtualTrans := e.applyUowShortcut(newComposition); virtualTrans != nil {
			transformations = append(transformations, virtualTrans)
		}
	}
	/**
	* Sometimes, a tone's position in a previous state must be changed to fit the new state
	*
	* e.g.
	* prev state: chuyr -> chuỷ
	* this state: chuyrene -> chuyển
	**/
	transformations = append(transformations, e.refreshLastToneTarget(append(composition, transformations...))...)
	return transformations
}

func (e *BambooEngine) newComposition(composition []*Transformation, key rune, isUpperCase bool) []*Transformation {
	// Just process the key stroke on the last syllable
	var previousTransformations, lastSyllable = extractLastSyllable(composition)

	// Find all possible transformations this keypress can generate
	lastSyllable = append(lastSyllable, e.generateTransformations(lastSyllable, key, isUpperCase)...)

	// Put these transformations back to the composition
	return append(previousTransformations, lastSyllable...)
}

func (e *BambooEngine) applyUowShortcut(syllable []*Transformation) *Transformation {
	str := Flatten(syllable, ToneLess|LowerCase)
	if len(e.inputMethod.SuperKeys) > 0 && regUOhTail.MatchString(str) {
		if target, missingRule := e.findTargetByKey(syllable, e.inputMethod.SuperKeys[0]); target != nil {
			missingRule.Key = rune(0) // virtual rule should not appear in the raw string
			virtualTrans := &Transformation{
				Rule:   missingRule,
				Target: target,
			}
			return virtualTrans
		}
	}
	return nil
}

func (e *BambooEngine) refreshLastToneTarget(syllable []*Transformation) []*Transformation {
	if e.flags&EfreeToneMarking != 0 && isValid(syllable, false) {
		return refreshLastToneTarget(syllable, e.flags&EstdToneStyle != 0)
	}
	return nil
}

/***** BEGIN SIDE-EFFECT METHODS ******/

func (e *BambooEngine) ProcessString(str string, mode Mode) {
	for _, key := range str {
		e.ProcessKey(key, mode)
	}
}

func (e *BambooEngine) ProcessKey(key rune, mode Mode) {
	var lowerKey = unicode.ToLower(key)
	var isUpperCase = unicode.IsUpper(key)
	if mode&EnglishMode != 0 || !e.CanProcessKey(lowerKey) {
		if mode&InReverseOrder != 0 {
			e.composition = append([]*Transformation{newAppendingTrans(lowerKey, isUpperCase)}, e.composition...)
			return
		}
		e.composition = append(e.composition, newAppendingTrans(lowerKey, isUpperCase))
		return
	}
	e.composition = e.newComposition(e.composition, lowerKey, isUpperCase)
}

func (e *BambooEngine) RestoreLastWord(toVietnamese bool) {
	var previous, lastComb = extractLastWord(e.composition, e.GetInputMethod().Keys)
	if len(lastComb) == 0 {
		return
	}
	if !toVietnamese {
		e.composition = append(previous, breakComposition(lastComb)...)
	} else {
		var newComp []*Transformation
		for _, tnx := range lastComb {
			newComp = e.newComposition(newComp, tnx.Rule.Key, tnx.IsUpperCase)
		}
		e.composition = append(previous, newComp...)
	}
}

func (e *BambooEngine) Reset() {
	e.composition = nil
}

// Find the last APPENDING transformation and all
// the transformations that add effects to it.
func (e *BambooEngine) RemoveLastChar(refreshLastToneTarget bool) {
	var lastAppending = findLastAppendingTrans(e.composition)
	if lastAppending == nil {
		return
	}
	if !e.CanProcessKey(lastAppending.Rule.Key) {
		e.composition = e.composition[:len(e.composition)-1]
		return
	}
	var previous, lastComb = extractLastWord(e.composition, e.GetInputMethod().Keys)
	var newComb []*Transformation
	for _, t := range lastComb {
		if t.Target == lastAppending || t == lastAppending {
			continue
		}
		newComb = append(newComb, t)
	}
	if refreshLastToneTarget {
		newComb = append(newComb, e.refreshLastToneTarget(newComb)...)
	}
	e.composition = append(previous, newComb...)
}

/***** END SIDE-EFFECT METHODS ******/
