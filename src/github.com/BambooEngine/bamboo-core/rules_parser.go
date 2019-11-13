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
	"XoaDauThanh": ToneNone,
	"DauSac":      ToneAcute,
	"DauHuyen":    ToneGrave,
	"DauNga":      ToneTilde,
	"DauNang":     ToneDot,
	"DauHoi":      ToneHook,
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
	MarkNone  Mark = iota << 0
	MarkHat   Mark = iota
	MarkBreve Mark = iota
	MarkHorn  Mark = iota
	MarkDash  Mark = iota
	MarkRaw   Mark = iota
)

type Tone uint8

const (
	ToneNone  Tone = iota << 0
	ToneGrave Tone = iota
	ToneAcute Tone = iota
	ToneHook  Tone = iota
	ToneTilde Tone = iota
	ToneDot   Tone = iota
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
	var tones = []Tone{ToneNone, ToneDot, ToneAcute, ToneGrave, ToneHook, ToneTilde}
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
