package main

import (
	"testing"

	"github.com/BambooEngine/bamboo-core"
)

type keyEvent struct {
	keys                [3]uint32
	canBeProcessed      bool
	expectedPreeditText string
	expectedCommitText  string
}

func asciiToKeys(s rune) [3]uint32 {
	return [3]uint32{uint32(s), uint32(s), 0}
}

func generateKeyEvents(s string, v []string, appendKeys ...keyEvent) []keyEvent {
	var kv []keyEvent
	for i, c := range s {
		kv = append(kv, keyEvent{keys: asciiToKeys(c), canBeProcessed: true, expectedPreeditText: v[i], expectedCommitText: v[i]})
	}
	return append(kv, appendKeys...)
}

func generateMetaKeyEvent(keys [3]uint32) func(expectedText ...string) keyEvent {
	return func(expectedText ...string) keyEvent {
		kv := keyEvent{keys: keys, canBeProcessed: false}
		if len(expectedText) > 0 {
			kv.expectedCommitText = expectedText[0]
		}
		if len(expectedText) > 1 {
			kv.expectedPreeditText = expectedText[1]
		}
		return kv
	}
}

var enter = generateMetaKeyEvent([3]uint32{0xff0d, 0xff0d, 0})
var control = generateMetaKeyEvent([3]uint32{0xffe3, 0xffe3, 0})

type testCase struct {
	name      string
	keyEvents []keyEvent
	inputMode int
	mTable    map[string]string
}

func TestPreeditEngine(t *testing.T) {
	for _, tc := range []testCase{
		{
			name: "empty_key_events",
			keyEvents: []keyEvent{
				{
					keys:           [3]uint32{0, 0, 0},
					canBeProcessed: false,
				},
			},
		},
		{
			name: "control_a",
			keyEvents: []keyEvent{
				control(),
				{keys: [3]uint32{0x0061, 0x0061, 4}, canBeProcessed: false, expectedPreeditText: ""},
			},
		},
		{
			name:   "macro_control_a",
			mTable: map[string]string{"->": "arrow"},
			keyEvents: []keyEvent{
				control(),
				{keys: [3]uint32{0x0061, 0x0061, 4}, canBeProcessed: false, expectedPreeditText: ""},
			},
		},
		{
			name:      "duowidro",
			keyEvents: generateKeyEvents("duowidro", []string{"d", "du", "duo", "dươ", "dươi", "đươi", "đưởi", "đuổi"}),
		},
		{
			name:      "duowidro_enter",
			keyEvents: generateKeyEvents("duowidro", []string{"d", "du", "duo", "dươ", "dươi", "đươi", "đưởi", "đuổi"}, enter("đuổi")),
		},
		{
			name:      "macro_vowl_space",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vowl ", []string{"v", "vo", "vơ", "vơl", ""}, control("vowl ")),
		},
		{
			name:      "macro_vowl_enter",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vowl", []string{"v", "vo", "vơ", "vơl"}, enter("vowl")),
		},
		{
			name:      "macro_duowidro_enter",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("duowidro", []string{"d", "du", "duo", "dươ", "dươi", "đươi", "đưởi", "đuổi"}, enter("đuổi")),
		},
		{
			name: "workaround_spreadsheet_number_enter",
			keyEvents: []keyEvent{
				{keys: asciiToKeys('1'), canBeProcessed: false, expectedPreeditText: ""},
				{keys: asciiToKeys('2'), canBeProcessed: false, expectedPreeditText: ""},
				enter(),
			},
		},
		{
			name:      "macro_vn_dot",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vn.", []string{"v", "vn", ""}, enter("việt nam.")),
		},
		{
			name:   "macro_vn_comma_space",
			mTable: map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vn", []string{"v", "vn"}, []keyEvent{
				{keys: asciiToKeys(','), canBeProcessed: true, expectedCommitText: "việt nam,"},
				{keys: asciiToKeys(' '), canBeProcessed: true, expectedCommitText: " "},
				enter("việt nam, "),
			}...),
		},
		{
			name:      "macro_vn_enter",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vn", []string{"v", "vn"}, enter("việt nam")),
		},
		{
			name:      "macro_arrow_dot",
			mTable:    map[string]string{"->": "arrow"},
			keyEvents: generateKeyEvents("->.", []string{"-", "->", ""}, control("arrow.")),
		},
		{
			name:      "macro_arrow_enter",
			mTable:    map[string]string{"->": "arrow"},
			keyEvents: generateKeyEvents("->", []string{"-", "->"}, enter("arrow")),
		},
		{
			name:   "macro_csao_space",
			mTable: map[string]string{"csao": "✪", "csao2": "✬"},
			keyEvents: generateKeyEvents("csao", []string{"c", "cs", "csa", "csao"}, []keyEvent{
				{keys: asciiToKeys(' '), canBeProcessed: true, expectedCommitText: "✪ "},
			}...),
		},
		{
			name:      "macro_csao2_enter",
			mTable:    map[string]string{"csao": "✪", "csao2": "✬"},
			keyEvents: generateKeyEvents("csao2", []string{"c", "cs", "csa", "csao", "csao2"}, enter("✬")),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.inputMode = preeditIM
			assertEngine(t, tc, func(t testing.TB, fe *fakeEngine, e IEngine) {
				for _, ev := range tc.keyEvents {
					keys := ev.keys
					t.Logf("Processing key %c %v", rune(keys[0]), keys)
					ret, _ := e.ProcessKeyEvent(keys[0], keys[1], keys[2])
					if ret != ev.canBeProcessed {
						t.Errorf("Is key can be processed? expected (%v), got (%v).", ev.canBeProcessed, ret)
					}
					if ev.canBeProcessed && fe.preeditText != ev.expectedPreeditText {
						t.Errorf("Preedit text, expected (%s), got (%s).", ev.expectedPreeditText, fe.preeditText)
					}
					if !ev.canBeProcessed && ev.expectedCommitText != fe.commitText {
						t.Errorf("Commit text, expected (%s), got (%s).", ev.expectedCommitText, fe.commitText)
					}
				}
			})
		})
	}
}

func TestBsEngine(t *testing.T) {
	for _, tc := range []testCase{
		{
			name: "empty_key_events",
			keyEvents: []keyEvent{
				{
					keys:           [3]uint32{0, 0, 0},
					canBeProcessed: false,
				},
			},
		},
		{
			name: "control_a",
			keyEvents: []keyEvent{
				control(),
				{keys: [3]uint32{0x0061, 0x0061, 4}, canBeProcessed: false}, // Ctrl+A
			},
		},
		{
			name:      "vn_dot_enter",
			keyEvents: generateKeyEvents("vn.", []string{"v", "vn", "vn."}, enter("vn.")),
		},
		{
			name:   "macro_control_a",
			mTable: map[string]string{"->": "arrow"},
			keyEvents: []keyEvent{
				control(),
				{keys: [3]uint32{0x0061, 0x0061, 4}, canBeProcessed: false}, // Ctrl+A
			},
		},
		{
			name:      "duowidro",
			keyEvents: generateKeyEvents("duowidro", []string{"d", "du", "duo", "dươ", "dươi", "đươi", "đưởi", "đuổi"}),
		},
		{
			name:      "duowidro_enter",
			keyEvents: generateKeyEvents("duowidro", []string{"d", "du", "duo", "dươ", "dươi", "đươi", "đưởi", "đuổi"}, enter("đuổi")),
		},
		{
			name:      "macro_vowl_space",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vowl ", []string{"v", "vo", "vơ", "vơl", "vowl "}),
		},
		{
			name:      "macro_vowl_enter",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vowl", []string{"v", "vo", "vơ", "vơl"}, enter("vowl")),
		},
		{
			name:      "macro_duowidro_enter",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("duowidro", []string{"d", "du", "duo", "dươ", "dươi", "đươi", "đưởi", "đuổi"}, enter("đuổi")),
		},
		{
			name: "workaround_spreadsheet_number_enter",
			keyEvents: []keyEvent{
				{keys: asciiToKeys('1'), canBeProcessed: false, expectedPreeditText: ""},
				{keys: asciiToKeys('2'), canBeProcessed: false, expectedPreeditText: ""},
				enter(),
			},
		},
		{
			name:   "macro_12",
			mTable: map[string]string{"vn": "việt nam"},
			keyEvents: []keyEvent{
				{keys: asciiToKeys('1'), canBeProcessed: true, expectedCommitText: "1"},
				{keys: asciiToKeys('2'), canBeProcessed: true, expectedCommitText: "12"},
				enter("12"),
			},
		},
		{
			name:      "macro_vn_dot",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vn.", []string{"v", "vn", "việt nam."}),
		},
		{
			name:      "macro_vn_comma_space",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vn, ", []string{"v", "vn", "việt nam,", "việt nam, "}),
		},
		{
			name:      "macro_vn_enter",
			mTable:    map[string]string{"vn": "việt nam"},
			keyEvents: generateKeyEvents("vn", []string{"v", "vn"}, enter("việt nam")),
		},
		{
			name:      "macro_arrow_dot",
			mTable:    map[string]string{"->": "arrow"},
			keyEvents: generateKeyEvents("->.", []string{"-", "->", "arrow."}),
		},
		{
			name:      "macro_arrow_enter",
			mTable:    map[string]string{"->": "arrow"},
			keyEvents: generateKeyEvents("->", []string{"-", "->"}, enter("arrow")),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.inputMode = surroundingTextIM
			assertEngine(t, tc, func(t testing.TB, fe *fakeEngine, e IEngine) {
				for _, ev := range tc.keyEvents {
					keys := ev.keys
					t.Logf("Processing key %c %v", rune(keys[0]), keys)
					ret, _ := e.ProcessKeyEvent(keys[0], keys[1], keys[2])
					if ret != ev.canBeProcessed {
						t.Errorf("Is key can be processed? expected (%v), got (%v).", ev.canBeProcessed, ret)
					}
					if fe.commitText != ev.expectedCommitText {
						t.Errorf("Commit text, expected (%s), got (%s).", ev.expectedCommitText, fe.commitText)
					}
				}
			})
		})
	}
}

func assertEngine(t testing.TB, tc testCase, assertFn func(testing.TB, *fakeEngine, IEngine)) {
	fe := NewFakeEngine()
	engineName := "test"
	var cfg = defaultCfg()
	cfg.DefaultInputMode = tc.inputMode
	inputMethod := bamboo.ParseInputMethod(cfg.InputMethodDefinitions, cfg.InputMethod)
	if tc.mTable != nil {
		cfg.IBflags |= IBmacroEnabled
	}
	e := NewIbusBambooEngine(engineName, &cfg, fe, bamboo.NewEngine(inputMethod, cfg.Flags))
	if tc.mTable != nil {
		e.macroTable = &MacroTable{
			mTable: tc.mTable,
		}
	}
	assertFn(t, fe, e)
}
