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
	EautoCorrectEnabled
	EddFreeStyle
	EstdFlags = EfreeToneMarking | EstdToneStyle | EautoCorrectEnabled | EddFreeStyle
)

type Transformation struct {
	Rule        Rule
	Target      *Transformation // For Tone/Mark transformation
	IsUpperCase bool
	IsDeleted   bool
	Dest        uint // For Appending, a pointer to the char in the flattened string made by this Trans
}

type IEngine interface {
	SetFlag(uint)
	GetInputMethod() InputMethod
	ProcessChar(rune, Mode)
	ProcessString(string, Mode)
	GetProcessedString(Mode, bool) string
	IsSpellingCorrect(Mode) bool
	GetSpellingMatchResult(Mode, bool) uint8
	HasTone() bool
	Reset()
	RemoveLastChar()
}

type BambooEngine struct {
	composition []*Transformation
	inputMethod InputMethod
	flags       uint
}

func NewEngine(im string, flag uint, dictionary map[string]bool) IEngine {
	inputMethod, found := InputMethods[im]
	if !found {
		panic("The input method is not supported")
	}
	engine := BambooEngine{
		inputMethod: inputMethod,
		flags:       flag,
	}
	for word, _ := range dictionary {
		AddTrie(spellingTrie, []rune(RemoveToneFromWord(word)), false)
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

func (e *BambooEngine) HasTone() bool {
	for _, t := range e.composition {
		if t.Rule.EffectType == ToneTransformation {
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

func (e *BambooEngine) getApplicableRules(key rune) []Rule {
	var applicableRules []Rule
	for _, inputRule := range e.inputMethod.Rules {
		if inputRule.Key == key {
			applicableRules = append(applicableRules, inputRule)
		}
	}
	return applicableRules
}

func (e *BambooEngine) findTargetForKey(key rune) (*Transformation, Rule) {
	var applicableRules = e.getApplicableRules(key)
	var lastAppending = findLastAppendingTrans(e.composition)
	for _, applicableRule := range applicableRules {
		var target *Transformation = nil
		if applicableRule.EffectType == MarkTransformation {
			return findMarkTarget(e.composition, applicableRules)
		} else if applicableRule.EffectType == ToneTransformation {
			if e.flags&EfreeToneMarking != 0 {
				if hasValidTone(e.composition, Tone(applicableRule.Effect)) {
					target = findToneTarget(e.composition, e.flags&EstdToneStyle != 0)
					if !isFree(e.composition, target, ToneTransformation) {
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

func (e *BambooEngine) GetSpellingMatchResult(mode Mode, deepSearch bool) uint8 {
	return getSpellingMatchResult(getLastWord(e.composition), mode, deepSearch)
}

func (e *BambooEngine) createCompositionForKey(chr rune, isUpperCase bool) []*Transformation {
	var transformations []*Transformation
	var appendingRule = findAppendingRule(e.inputMethod.Rules, chr)
	if unicode.IsUpper(appendingRule.EffectOn) {
		isUpperCase = true
		appendingRule.EffectOn = unicode.ToLower(appendingRule.EffectOn)
	}
	transformations = e.createCompositionForRule(appendingRule, isUpperCase)
	return transformations
}

func (e *BambooEngine) GetRawString() string {
	var seq []rune
	for _, t := range e.composition {
		seq = append(seq, t.Rule.Key)
	}
	return string(seq)
}

func (e *BambooEngine) GetProcessedString(mode Mode, lastWordOnly bool) string {
	if lastWordOnly {
		var lastComb = getLastWord(e.composition)
		if len(lastComb) > 0 {
			return Flatten(lastComb, mode)
		}
		return ""
	}
	return Flatten(e.composition, mode)
}

func (e *BambooEngine) IsSupportedKey(key rune) bool {
	if IsAlpha(key) {
		return true
	}
	if inKeyMap(e.GetInputMethod().Keys, key) {
		return true
	}
	return false
}

/***** BEGIN SIDE-EFFECT METHODS ******/

func (e *BambooEngine) ProcessChar(key rune, mode Mode) {
	var isUpperCase bool
	if unicode.IsUpper(key) {
		isUpperCase = true
	}
	key = unicode.ToLower(key)
	if mode&EnglishMode != 0 || !e.IsSupportedKey(key) {
		e.composition = append(e.composition, createAppendingTrans(key, isUpperCase))
		return
	}
	var previousComposition []*Transformation
	if len(e.composition) > 0 {
		var lastComb = getLastCombination(getLastWord(e.composition))
		if len(lastComb) > 0 {
			lastComb = getLastCombination(getLastWord(e.composition))
			var idx = findTransformationIndex(e.composition, lastComb[0])
			if idx > 0 {
				previousComposition = e.composition[:idx]
				e.composition = lastComb
			}
		} else {
			previousComposition = e.composition
			e.composition = nil
		}
	}
	if len(e.composition) > 0 && e.isEffectiveKey(key) {
		// garbage collection
		e.composition = freeComposition(e.composition)

		if target, _ := e.findTargetForKey(key); target == nil {
			if key == e.composition[len(e.composition)-1].Rule.Key {
				// Double typing an effect key undoes it and its effects.
				e.composition = undoesTransformations(e.composition, e.getApplicableRules(key))
				e.composition = append(e.composition, createAppendingTrans(key, isUpperCase))

				if previousComposition != nil {
					e.composition = append(previousComposition, e.composition...)
				}
				return
			} else {
				// Or an effect key may override other effect keys
				e.composition = undoesTransformations(e.composition, e.getApplicableRules(key))
			}
		}
	}
	// TODO: need to refactor
	if e.flags&EautoCorrectEnabled != 0 && (e.isSuperKey(key) || (!e.isToneKey(key) && hasSuperWord(e.composition))) {
		if missingRule, found := findMissingRuleForUo(e.composition, e.isSuperKey(key)); found {
			var targets = findMarkTargets(e.composition, missingRule)
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
	transformations := e.createCompositionForKey(key, isUpperCase)
	e.composition = append(e.composition, transformations...)

	/**
	* Sometimes, a tone's position in a previous state must be changed to fit the new state
	*
	* e.g.
	* prev state: chuyenr -> chuỷen
	* this state: chuyenre -> chuyển
	**/
	if e.flags&EstdToneStyle != 0 && shouldRefreshLastToneTarget(e.composition) {
		e.composition = refreshLastToneTarget(e.composition)
	}

	if previousComposition != nil {
		e.composition = append(previousComposition, e.composition...)
	}
}

func (e *BambooEngine) ProcessString(str string, mode Mode) {
	for _, chr := range []rune(str) {
		e.ProcessChar(chr, mode)
	}
}

func (e *BambooEngine) Reset() {
	e.composition = []*Transformation{}
}

func (e *BambooEngine) RemoveLastChar() {
	var lastAppending = findLastAppendingTrans(e.composition)
	var transformations = getTransformationsTargetTo(e.composition, lastAppending)
	for _, trans := range append(transformations, lastAppending) {
		e.composition = removeTrans(e.composition, trans)
	}
}

/***** END SIDE-EFFECT METHODS ******/
