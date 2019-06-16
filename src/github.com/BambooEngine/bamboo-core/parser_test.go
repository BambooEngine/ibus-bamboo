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
	"testing"
)

func TestParseToneRules(t *testing.T) {
	rules := ParseRules('z', "XoaDauThanh")
	if len(rules) != 1 || rules[0].EffectType != ToneTransformation || Tone(rules[0].Effect) != TONE_NONE {
		t.Errorf("Test parse None Rule. Got %v, expected %v", rules[0], Rule{
			Key:        'z',
			EffectType: ToneTransformation,
			Effect:     0,
		})
	}
	rules = ParseRules('x', "DauNga")
	if len(rules) != 1 || rules[0].EffectType != ToneTransformation || rules[0].GetTone() != TONE_TILDE {
		t.Errorf("Test parse None Rule. Got %v, expected %v", rules[0], Rule{
			Key:        'x',
			EffectType: ToneTransformation,
			Effect:     uint8(TONE_TILDE),
		})
	}
}

func TestParseTonelessRules(t *testing.T) {
	rules := ParseTonelessRules('d', "D_Đ")
	idx := 0
	if len(rules) != 2 || rules[idx].EffectType != MarkTransformation || rules[idx].Effect != uint8(MARK_DASH) || rules[idx].EffectOn != 'd' {
		t.Errorf("Test parsing Mark Rule. Got %v, expected %v", rules[idx], Rule{
			Key:        'd',
			EffectType: MarkTransformation,
			Effect:     uint8(MARK_DASH),
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
	t.Log("RULES=", rules)
	if len(rules) != 33 {
		t.Errorf("Test the length of parsing mark rule. Got %d, expected %d", len(rules), 30)
	}
	if rules[0].EffectType != MarkTransformation || rules[0].GetMark() != MARK_HORN || rules[0].EffectOn != 'u' {
		t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[0], Rule{
			Key:        'w',
			EffectType: MarkTransformation,
			Effect:     uint8(MARK_HORN),
			EffectOn:   'u',
		})
	}
	idx = 7
	if rules[idx].EffectType != MarkTransformation || rules[idx].GetMark() != MARK_HORN || rules[idx].EffectOn != 'o' {
		t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[idx], Rule{
			Key:        'w',
			EffectType: MarkTransformation,
			Effect:     uint8(MARK_HORN),
			EffectOn:   'o',
		})
	}
	idx = 20
	if rules[idx].EffectType != MarkTransformation || rules[idx].GetMark() != MARK_BREVE || rules[idx].EffectOn != 'a' {
		t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[idx], Rule{
			Key:        'w',
			EffectType: MarkTransformation,
			Effect:     uint8(MARK_BREVE),
			EffectOn:   'a',
		})
	}
	rules = ParseTonelessRules('w', "UOA_ƯƠĂ__Ư")
	if len(rules) != 34 {
		t.Errorf("Test the length of parsing mark rule. Got %d, expected %d", len(rules), 31)
	} else {
		t.Log("RULES[UOA_ƯƠĂ__Ư]=", rules)
		idx = 20
		if rules[idx].EffectType != MarkTransformation || rules[idx].GetMark() != MARK_BREVE || rules[idx].EffectOn != 'a' {
			t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[idx], Rule{
				Key:        'w',
				EffectType: MarkTransformation,
				Effect:     uint8(MARK_BREVE),
				EffectOn:   'a',
			})
		}
		idx = 33
		if rules[idx].EffectType != Appending || rules[idx].EffectOn != 'ư' {
			t.Errorf("Test parsing mark Rule. Got %v, expected %v", rules[idx], Rule{
				Key:        'w',
				EffectType: Appending,
				EffectOn:   'ư',
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
