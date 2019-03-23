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
	"unicode"
)

type Mode uint

const (
	VietnameseMode Mode = 1 << iota
	EnglishMode
	ToneLess
	MarkLess
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
	EstdFlags = EfreeToneMarking | EstdToneStyle | EautoCorrectEnabled
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
	ProcessKey(rune, Mode)
	ProcessString(string, Mode)
	GetProcessedString(Mode, bool) string
	GetSpellingMatchResult(Mode, bool) uint8
	CanProcessKey(rune) bool
	HasTone() bool
	Reset()
	RemoveLastChar()
	GetRawString() string
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
	for word := range dictionary {
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

func (e *BambooEngine) HasTone() bool {
	for _, t := range getLastWord(e.composition, nil) {
		if t.Rule.EffectType == ToneTransformation {
			return true
		}
	}
	return false
}

func (e *BambooEngine) isSuperKey(key rune) bool {
	return inKeyList(e.GetInputMethod().SuperKeys, key)
}

func (e *BambooEngine) isSupportedKey(key rune) bool {
	if IsAlpha(key) {
		return true
	}
	return inKeyList(e.GetInputMethod().Keys, key)
}

func (e *BambooEngine) isToneKey(key rune) bool {
	return inKeyList(e.GetInputMethod().ToneKeys, key)
}

func (e *BambooEngine) isEffectiveKey(key rune) bool {
	return inKeyList(e.GetInputMethod().Keys, key)
}

func (e *BambooEngine) GetSpellingMatchResult(mode Mode, deepSearch bool) uint8 {
	return getSpellingMatchResult(getLastWord(e.composition, nil), mode, deepSearch)
}

func (e *BambooEngine) GetRawString() string {
	var seq []rune
	for _, t := range e.composition {
		seq = append(seq, t.Rule.Key)
	}
	return string(seq)
}

func (e *BambooEngine) GetProcessedString(mode Mode, lastWordOnly bool) string {
	var effectiveKeys = e.inputMethod.Keys
	if lastWordOnly {
		effectiveKeys = nil
	}
	var lastComb = getLastWord(e.composition, effectiveKeys)
	if len(lastComb) > 0 {
		return Flatten(lastComb, mode)
	}
	return ""
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

func (e *BambooEngine) findTargetFromKey(composition []*Transformation, key rune) (*Transformation, Rule) {
	return findTargetFromKey(composition, e.getApplicableRules(key), e.flags)
}

func (e *BambooEngine) createAppendingTransFromKey(key rune, isUpperCase bool) *Transformation {
	var appendingRule = findAppendingRule(e.inputMethod.Rules, key)
	if unicode.IsUpper(appendingRule.EffectOn) {
		isUpperCase = true
		appendingRule.EffectOn = unicode.ToLower(appendingRule.EffectOn)
	}
	var trans = new(Transformation)
	trans.Rule = appendingRule
	trans.IsUpperCase = isUpperCase
	return trans
}

// Find all possible transformations this keypress can generate
func (e *BambooEngine) createTransformations(composition []*Transformation, key rune, isUpperCase bool) []*Transformation {
	return generateTransformations(composition, e.getApplicableRules(key), e.createAppendingTransFromKey(key, isUpperCase), e.flags)
}

func (e *BambooEngine) isTransformationForUoMissed(composition []*Transformation) bool {
	return e.flags&EautoCorrectEnabled != 0 &&
		len(composition) > 0 &&
		hasSuperWord(composition) &&
		getSpellingMatchResult(composition, ToneLess, false) == FindResultMatchPrefix
}

func (e *BambooEngine) CanProcessKey(key rune) bool {
	return e.isSupportedKey(key)
}

/***** BEGIN SIDE-EFFECT METHODS ******/

func (e *BambooEngine) ProcessKey(key rune, mode Mode) {
	var isUpperCase bool
	if unicode.IsUpper(key) {
		isUpperCase = true
	}
	key = unicode.ToLower(key)
	if mode&EnglishMode != 0 || !e.isSupportedKey(key) {
		e.composition = append(e.composition, createAppendingTrans(key, isUpperCase))
		return
	}
	// Extract the last syllable from the composition
	var lastSyllable, previousTransformations = extractLastSyllable(e.composition)

	// Implement the double typing an effective key
	if len(lastSyllable) > 0 && e.isEffectiveKey(key) {
		// remove unused transformations
		lastSyllable = freeComposition(lastSyllable)

		if target, _ := e.findTargetFromKey(lastSyllable, key); target == nil {
			if key == lastSyllable[len(lastSyllable)-1].Rule.Key {
				// Double typing an effect key undoes it and its effects.
				lastSyllable = undoesTransformations(lastSyllable, e.getApplicableRules(key))
				lastSyllable = append(lastSyllable, createAppendingTrans(key, isUpperCase))

				e.composition = append(previousTransformations, lastSyllable...)
				return
			} else {
				// Or an effect key may override other effect keys
				lastSyllable = undoesTransformations(lastSyllable, e.getApplicableRules(key))
			}
		}
	}

	// Just process the key stroke on the last syllable
	lastSyllable = append(lastSyllable, e.createTransformations(lastSyllable, key, isUpperCase)...)

	// Implement the uow typing shortcut by creating a virtual
	// Mark.HORN rule that targets 'u' or 'o'.
	if e.isTransformationForUoMissed(lastSyllable) {
		if target, missingRule := e.findTargetFromKey(lastSyllable, e.inputMethod.SuperKeys[0]); target != nil {
			missingRule.Key = rune(0) // virtual rule should not appear in the raw string
			virtualTrans := &Transformation{
				Rule:   missingRule,
				Target: target,
			}
			lastSyllable = append(lastSyllable, virtualTrans)
		}
	}
	/**
	* Sometimes, a tone's position in a previous state must be changed to fit the new state
	*
	* e.g.
	* prev state: chuyenr -> chuỷen
	* this state: chuyenre -> chuyển
	**/
	if e.flags&EstdToneStyle != 0 && shouldRefreshLastToneTarget(lastSyllable) {
		lastSyllable = refreshLastToneTarget(lastSyllable)
	}

	e.composition = append(previousTransformations, lastSyllable...)
}

func (e *BambooEngine) ProcessString(str string, mode Mode) {
	for _, key := range []rune(str) {
		e.ProcessKey(key, mode)
	}
}

func (e *BambooEngine) Reset() {
	e.composition = nil
}

// Find the last APPENDING transformation and all
// the transformations that add effects to it.
func (e *BambooEngine) RemoveLastChar() {
	var lastAppending = findLastAppendingTrans(e.composition)
	if lastAppending == nil {
		return
	}
	var transformations = getTransformationsTargetTo(e.composition, lastAppending)
	for _, trans := range append(transformations, lastAppending) {
		e.composition = removeTrans(e.composition, trans)
	}
}

/***** END SIDE-EFFECT METHODS ******/
