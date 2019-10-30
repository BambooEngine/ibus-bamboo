/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENSE>.
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
	Replacing          EffectType = iota
)

// type alias
type Mark uint8

const (
	MARK_NONE  Mark = iota << 0
	MARK_HAT   Mark = iota
	MARK_BREVE Mark = iota
	MARK_HORN  Mark = iota
	MARK_DASH  Mark = iota
	MARK_RAW   Mark = iota
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
	Result        rune
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
			rules = append(rules, ParseToneLessRule(key, effectiveOn, results[i], effect)...)
		}
		if rule, ok := getAppendingRule(key, parts[3]); ok {
			rules = append(rules, rule)
		}

	} else if rule, ok := getAppendingRule(key, line); ok {
		rules = append(rules, rule)
	}
	return rules
}

func ParseToneLessRule(key, effectiveOn, result rune, effect Mark) []Rule {
	var rules []Rule
	var tones = []Tone{TONE_NONE, TONE_DOT, TONE_ACUTE, TONE_GRAVE, TONE_HOOK, TONE_TILDE}
	for _, chr := range getMarkFamily(effectiveOn) {
		if chr == result {
			var rule Rule
			rule.Key = key
			rule.EffectType = MarkTransformation
			rule.Effect = 0
			rule.EffectOn = result
			rule.Result = effectiveOn
			rules = append(rules, rule)
		} else if IsVowel(chr) {
			for tone := range tones {
				var rule Rule
				rule.Key = key
				rule.EffectType = MarkTransformation
				rule.EffectOn = AddToneToChar(chr, uint8(tone))
				rule.Effect = uint8(effect)
				rule.Result = AddToneToChar(result, uint8(tone))
				rules = append(rules, rule)
			}
		} else {
			var rule Rule
			rule.Key = key
			rule.EffectType = MarkTransformation
			rule.EffectOn = chr
			rule.Effect = uint8(effect)
			rule.Result = result
			rules = append(rules, rule)
		}
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
		rule.Result = chars[0]
		if len(chars) > 1 {
			for _, chr := range chars[1:] {
				rule.AppendedRules = append(rule.AppendedRules, Rule{
					Key:        key,
					EffectType: Appending,
					EffectOn:   chr,
					Result:     chr,
				})
			}
		}
		return rule, true
	}
	return rule, false
}
