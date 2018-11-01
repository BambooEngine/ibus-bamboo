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
	"testing"
)

func TestParseToneRules(t *testing.T) {
	rules := ParseRules('z', "XoaDauThanh")
	if len(rules) != 1 || rules[0].EffectType != ToneTransformation || rules[0].Effect != TONE_NONE {
		t.Errorf("Test parse None Rule. Got %v, expected %v", rules[0], Rule{
			Key:        'z',
			EffectType: ToneTransformation,
			Effect:     TONE_NONE,
		})
	}
	rules = ParseRules('x', "DauNga")
	if len(rules) != 1 || rules[0].EffectType != ToneTransformation || rules[0].Effect != TONE_TILDE {
		t.Errorf("Test parse None Rule. Got %v, expected %v", rules[0], Rule{
			Key:        'x',
			EffectType: ToneTransformation,
			Effect:     TONE_TILDE,
		})
	}
}

func TestParseTonelessRules(t *testing.T) {
	rules := ParseTonelessRules('d', "D_Đ")
	if len(rules) != 1 || rules[0].EffectType != MarkTransformation || rules[0].Effect != MARK_DASH || rules[0].EffectOn != 'd' {
		t.Errorf("Test parsing Mark Rule. Got %v, expected %v", rules, Rule{
			Key:        'd',
			EffectType: MarkTransformation,
			Effect:     MARK_DASH,
			EffectOn:   'd',
		})
	}
	rules = ParseTonelessRules('{', "_Ư")
	if len(rules) != 1 || rules[0].EffectType != Appending || rules[0].EffectOn != 'Ư' {
		t.Errorf("Test parsing Append Rule. Got %v, expected %v", rules, Rule{
			Key:        '{',
			EffectType: Appending,
			EffectOn:   'Ư',
		})
	}
	rules = ParseTonelessRules('w', "UOA_ƯƠĂ")
	if len(rules) != 3 {
		t.Errorf("Test the length of parsing mark rule. Got %d, expected %d", len(rules), 3)
	}
	if rules[0].EffectType != MarkTransformation || rules[0].Effect != MARK_HORN || rules[0].EffectOn != 'u' {
		t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[0], Rule{
			Key:        'w',
			EffectType: MarkTransformation,
			Effect:     MARK_HORN,
			EffectOn:   'u',
		})
	}
	if rules[1].EffectType != MarkTransformation || rules[1].Effect != MARK_HORN || rules[1].EffectOn != 'o' {
		t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[1], Rule{
			Key:        'w',
			EffectType: MarkTransformation,
			Effect:     MARK_HORN,
			EffectOn:   'o',
		})
	}
	if rules[2].EffectType != MarkTransformation || rules[2].Effect != MARK_BREVE || rules[2].EffectOn != 'a' {
		t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[2], Rule{
			Key:        'w',
			EffectType: MarkTransformation,
			Effect:     MARK_BREVE,
			EffectOn:   'a',
		})
	}
	rules = ParseTonelessRules('w', "UOA_ƯƠĂ__Ư")
	if len(rules) != 4 {
		t.Errorf("Test the length of parsing mark rule. Got %d, expected %d", len(rules), 4)
	} else {
		if rules[2].EffectType != MarkTransformation || rules[2].Effect != MARK_BREVE || rules[2].EffectOn != 'a' {
			t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[2], Rule{
				Key:        'w',
				EffectType: MarkTransformation,
				Effect:     MARK_BREVE,
				EffectOn:   'a',
			})
		}
		if rules[3].EffectType != Appending || rules[3].EffectOn != 'ư' {
			t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[3], Rule{
				Key:        'w',
				EffectType: Appending,
			})
		}
	}

}

func TestAppendRule(t *testing.T) {
	rules := ParseTonelessRules('[', "__ươ")
	if len(rules) != 1 {
		t.Errorf("Test the length of parsing mark rule. Got %d, expected %d", len(rules), 1)
	} else {
		appendRules := rules[0].AppendedRules
		if len(appendRules) != 1 || appendRules[0].EffectType != Appending || appendRules[0].EffectOn != 'ơ' {
			t.Errorf("Test parsing append mark Rule. Got %v, expected %v", appendRules, Rule{
				Key:        '[',
				EffectType: Appending,
				EffectOn:   'ơ',
			})
		}
	}

	rules = ParseTonelessRules('{', "__ƯƠ")
	if len(rules) != 1 {
		t.Errorf("Test the length of parsing mark rule. Got %d, expected %d", len(rules), 1)
	} else {
		appendRules := rules[0].AppendedRules
		if len(appendRules) != 1 || appendRules[0].EffectType != Appending || appendRules[0].EffectOn != 'Ơ' {
			t.Errorf("Test parsing append mark Rule. Got %v, expected %v", appendRules, Rule{
				Key:        '{',
				EffectType: Appending,
				EffectOn:   'Ơ',
			})
		}
	}
}

func TestParseRulesWithIm(t *testing.T) {
}
