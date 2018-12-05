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

var InputMethods map[string]InputMethod

func init() {
	InputMethods = make(map[string]InputMethod, len(inputMethodDefinitions))
	for name, imDefinition := range inputMethodDefinitions {
		var im InputMethod
		im.Name = name
		for key, line := range imDefinition {
			im.Rules = append(im.Rules, ParseRules(key, line)...)
			if strings.Contains(strings.ToLower(line), "uo") {
				im.SuperKeys = append(im.SuperKeys, key)
			}
			if _, ok := tones[line]; ok {
				im.ToneKeys = append(im.ToneKeys, key)
			}
			im.Keys = append(im.Keys, key)
		}
		InputMethods[name] = im
	}
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

func ParseTonelessRules(key rune, line string) []Rule {
	var rules []Rule
	reg := regexp.MustCompile(`([a-zA-Z]+)_(\p{L}+)([_\p{L}]*)`)
	if reg.MatchString(line) {
		parts := reg.FindStringSubmatch(strings.ToLower(line))
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

func getAppendingRule(key rune, value string) (Rule, bool) {
	var rule Rule
	reg := regexp.MustCompile(`(_?)_(\p{L}+)`)
	if reg.MatchString(value) {
		parts := reg.FindStringSubmatch(value)
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
