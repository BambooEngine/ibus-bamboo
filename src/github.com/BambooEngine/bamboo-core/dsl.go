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
)

var tones = map[string]Tone{
	"XoaDauThanh": TONE_NONE,
	"DauSac":      TONE_ACUTE,
	"DauHuyen":    TONE_GRAVE,
	"DauNga":      TONE_TILDE,
	"DauNang":     TONE_DOT,
	"DauHoi":      TONE_HOOK,
}

type EffectType int

const (
	Appending          EffectType = iota << 0
	MarkTransformation EffectType = iota
	ToneTransformation EffectType = iota
)

// type alias
type Mark uint8

const (
	MARK_NONE  Mark = iota << 0
	MARK_HAT   Mark = iota
	MARK_BREVE Mark = iota
	MARK_HORN  Mark = iota
	MARK_DASH  Mark = iota
)

type Tone uint8

const (
	TONE_NONE  Tone = iota << 0
	TONE_GRAVE Tone = iota
	TONE_ACUTE Tone = iota
	TONE_HOOK  Tone = iota
	TONE_TILDE Tone = iota
	TONE_DOT   Tone = iota
)

type Rule struct {
	Key           rune
	Effect        uint8 // (Tone, Mark)
	EffectType    EffectType
	EffectOn      rune
	AppendedRules []Rule
}

func (r *Rule) SetTone(tone Tone) {
	r.Effect = uint8(tone)
}

func (r *Rule) SetMark(mark Mark) {
	r.Effect = uint8(mark)
}

func (r *Rule) GetTone() Tone {
	return Tone(r.Effect)
}

func (r *Rule) GetMark() Mark {
	return Mark(r.Effect)
}

type InputMethod struct {
	Name      string
	Rules     []Rule
	SuperKeys []rune
	ToneKeys  []rune
	Keys      []rune
}

func ParseInputMethod(imDef map[string]InputMethodDefinition, imName string) InputMethod {
	var inputMethods = parseInputMethods(imDef)
	if inputMethod, found := inputMethods[imName]; found {
		return inputMethod
	}
	return InputMethod{}
}

func parseInputMethods(imDef map[string]InputMethodDefinition) map[string]InputMethod {
	var inputMethods = make(map[string]InputMethod, len(imDef))
	for name, imDefinition := range imDef {
		var im InputMethod
		im.Name = name
		for keyStr, line := range imDefinition {
			var keys = []rune(keyStr)
			if len(keys) == 0 {
				continue
			}
			var key = keys[0]
			im.Rules = append(im.Rules, ParseRules(key, line)...)
			if strings.Contains(strings.ToLower(line), "uo") {
				im.SuperKeys = append(im.SuperKeys, key)
			}
			if _, ok := tones[line]; ok {
				im.ToneKeys = append(im.ToneKeys, key)
			}
			im.Keys = append(im.Keys, key)
		}
		inputMethods[name] = im
	}
	return inputMethods
}

func ParseRules(key rune, line string) []Rule {
	var rules []Rule
	if tone, ok := tones[line]; ok {
		var rule Rule
		rule.Key = key
		rule.EffectType = ToneTransformation
		rule.Effect = uint8(tone)
		rules = append(rules, rule)
	} else {
		rules = ParseTonelessRules(key, line)
	}
	return rules
}

var regDsl = regexp.MustCompile(`([a-zA-Z]+)_(\p{L}+)([_\p{L}]*)`)

func ParseTonelessRules(key rune, line string) []Rule {
	var rules []Rule
	if regDsl.MatchString(line) {
		parts := regDsl.FindStringSubmatch(strings.ToLower(line))
		effectiveOns := []rune(parts[1])
		results := []rune(parts[2])
		for i, effectiveOn := range effectiveOns {
			effect, found := FindMarkFromChar(results[i])
			if !found {
				continue
			}
			var rule Rule
			rule.Key = key
			rule.EffectType = MarkTransformation
			rule.EffectOn = effectiveOn
			rule.Effect = uint8(effect)

			rules = append(rules, rule)
		}
		if rule, ok := getAppendingRule(key, parts[3]); ok {
			rules = append(rules, rule)
		}

	} else if rule, ok := getAppendingRule(key, line); ok {
		rules = append(rules, rule)
	}
	return rules
}

var regDslAppending = regexp.MustCompile(`(_?)_(\p{L}+)`)

func getAppendingRule(key rune, value string) (Rule, bool) {
	var rule Rule
	if regDslAppending.MatchString(value) {
		parts := regDslAppending.FindStringSubmatch(value)
		chars := []rune(parts[2])
		rule.Key = key
		rule.EffectType = Appending
		rule.EffectOn = chars[0]
		if len(chars) > 1 {
			for _, chr := range chars[1:] {
				rule.AppendedRules = append(rule.AppendedRules, Rule{
					Key:        key,
					EffectType: Appending,
					EffectOn:   chr,
				})
			}
		}
		return rule, true
	}
	return rule, false
}
