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

import "unicode"

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
	Target      *Transformation
	IsUpperCase bool
}

type IEngine interface {
	SetFlag(uint)
	GetInputMethod() InputMethod
	ProcessKey(rune, Mode)
	ProcessString(string, Mode)
	GetProcessedString(Mode, bool) string
	GetSpellingMatchResult(Mode, bool) uint8
	CanProcessKey(rune) bool
	RemoveLastChar()
	RestoreLastWord()
	GetRawString() string
	Reset()
}

type BambooEngine struct {
	composition []*Transformation
	inputMethod InputMethod
	flags       uint
}

func NewEngine(inputMethod InputMethod, flag uint, dictionary map[string]bool) IEngine {
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
	return getSpellingMatchResult(getLastWord(e.composition, e.inputMethod.Keys), mode, deepSearch)
}

func (e *BambooEngine) GetRawString() string {
	var seq []rune
	for _, t := range e.composition {
		seq = append(seq, t.Rule.Key)
	}
	return string(seq)
}

func (e *BambooEngine) GetProcessedString(mode Mode, letterOnly bool) string {
	var effectiveKeys = e.inputMethod.Keys
	if letterOnly {
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
		if inputRule.Key == unicode.ToLower(key) {
			applicableRules = append(applicableRules, inputRule)
		}
	}
	return applicableRules
}

func (e *BambooEngine) findTargetByKey(composition []*Transformation, key rune) (*Transformation, Rule) {
	return findTarget(composition, e.getApplicableRules(key), e.flags)
}

// Find all possible transformations this keypress can generate
func (e *BambooEngine) generateTransformations(composition []*Transformation, lowerKey rune, isUpperCase bool) []*Transformation {
	return generateTransformations(composition, e.getApplicableRules(lowerKey), e.flags, lowerKey, isUpperCase)
}

func (e *BambooEngine) CanProcessKey(key rune) bool {
	return e.isSupportedKey(key)
}

/***** BEGIN SIDE-EFFECT METHODS ******/

func (e *BambooEngine) ProcessKey(key rune, mode Mode) {
	var lowerKey = unicode.ToLower(key)
	var isUpperCase = unicode.IsUpper(key)
	if mode&EnglishMode != 0 || !e.isSupportedKey(lowerKey) {
		e.composition = append(e.composition, newAppendingTrans(lowerKey, isUpperCase))
		return
	}
	// Just process the key stroke on the last syllable
	var lastSyllable, previousTransformations = extractLastSyllable(e.composition)

	lastSyllable = append(lastSyllable, e.generateTransformations(lastSyllable, lowerKey, isUpperCase)...)
	lastSyllable = e.refreshLastToneTarget(e.applyUowShortcut(lastSyllable))
	e.composition = append(previousTransformations, lastSyllable...)
}

// Implement the uow typing shortcut by creating a virtual
// Mark.HORN rule that targets 'u' or 'o'.
func (e *BambooEngine) applyUowShortcut(syllable []*Transformation) []*Transformation {
	if e.flags&EautoCorrectEnabled != 0 && len(e.inputMethod.SuperKeys) > 0 && isTransformationForUoMissed(syllable) {
		if target, missingRule := e.findTargetByKey(syllable, e.inputMethod.SuperKeys[0]); target != nil {
			missingRule.Key = rune(0) // virtual rule should not appear in the raw string
			virtualTrans := &Transformation{
				Rule:   missingRule,
				Target: target,
			}
			syllable = append(syllable, virtualTrans)
		}
	}
	return syllable
}

/**
* Sometimes, a tone's position in a previous state must be changed to fit the new state
*
* e.g.
* prev state: chuyr -> chuỷ
* this state: chuyrene -> chuyển
**/
func (e *BambooEngine) refreshLastToneTarget(syllable []*Transformation) []*Transformation {
	if e.flags&EfreeToneMarking != 0 {
		syllable = refreshLastToneTarget(syllable, e.flags&EstdToneStyle != 0)
	}
	return syllable
}

func (e *BambooEngine) ProcessString(str string, mode Mode) {
	for _, key := range []rune(str) {
		e.ProcessKey(key, mode)
	}
}

func (e *BambooEngine) RestoreLastWord() {
	var lastComb, previous = extractLastWord(e.composition, e.inputMethod.Keys)
	if len(lastComb) == 0 {
		return
	}
	e.composition = append(previous, breakComposition(lastComb)...)
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
	if !e.CanProcessKey(lastAppending.Rule.Key) {
		e.composition = e.composition[:len(e.composition)-1]
		return
	}
	var lastComb, previous = extractLastWord(e.composition, e.inputMethod.Keys)
	var newComb []*Transformation
	for _, t := range lastComb {
		if t.Target == lastAppending || t == lastAppending {
			continue
		}
		newComb = append(newComb, t)
	}
	newComb = e.refreshLastToneTarget(newComb)
	e.composition = append(previous, newComb...)
}

/***** END SIDE-EFFECT METHODS ******/
