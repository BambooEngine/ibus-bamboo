package main

import (
	"ibus-bamboo/config"
	"os"
	"sync/atomic"
	"testing"

	"github.com/BambooEngine/bamboo-core"
	ibus "github.com/BambooEngine/goibus"
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
			tc.inputMode = config.PreeditIM
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
			tc.inputMode = config.SurroundingTextIM
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

func TestInputModeSwitchWithoutWmClass(t *testing.T) {
	atomic.StoreInt32(&sessionInputMode, 0)
	defer atomic.StoreInt32(&sessionInputMode, 0)

	e := newTestEngine(config.PreeditIM)
	if e.getWmClass() != "" {
		t.Fatalf("wm class = %q, want empty", e.getWmClass())
	}

	handled, err := e.ProcessKeyEvent('~', '~', IBusShiftMask)
	if err != nil || !handled || !e.isInputModeLTOpened || e.inputModeLookupTable == nil {
		t.Fatalf("Shift+~ should open the mode table without a wm class: handled=%v opened=%v err=%v", handled, e.isInputModeLTOpened, err)
	}
	if len(e.inputModeLookupTable.Candidates) != len(config.ImLookupTable) {
		t.Fatalf("candidate count = %d, want %d", len(e.inputModeLookupTable.Candidates), len(config.ImLookupTable))
	}
	us := e.inputModeLookupTable.Candidates[config.UsIM-1].Value().(ibus.Text)
	if us.Text != config.ImLookupTable[config.UsIM] {
		t.Fatalf("exclusion label = %q, want %q", us.Text, config.ImLookupTable[config.UsIM])
	}

	handled, err = e.ProcessKeyEvent('2', '2', 0)
	if err != nil || !handled || e.isInputModeLTOpened {
		t.Fatalf("selecting a mode failed: handled=%v opened=%v err=%v", handled, e.isInputModeLTOpened, err)
	}
	if _, saved := e.config.InputModeMapping[""]; saved || len(e.config.InputModeMapping) != 0 {
		t.Fatalf("per-application map = %#v, want empty", e.config.InputModeMapping)
	}
	if got := e.getInputMode(); got != config.SurroundingTextIM {
		t.Fatalf("mode after the table closed = %d, want %d", got, config.SurroundingTextIM)
	}

	other := newTestEngine(config.PreeditIM)
	if got := other.getInputMode(); got != config.SurroundingTextIM {
		t.Fatalf("session mode on a later engine = %d, want %d", got, config.SurroundingTextIM)
	}
}

func TestExclusionListTogglesOff(t *testing.T) {
	atomic.StoreInt32(&sessionInputMode, 0)
	defer atomic.StoreInt32(&sessionInputMode, 0)

	e := newTestEngine(config.PreeditIM)
	openModeTable(t, e)
	selectMode(t, e, '7')
	if got := e.getInputMode(); got != config.UsIM {
		t.Fatalf("mode after adding an exclusion = %d, want %d", got, config.UsIM)
	}

	openModeTable(t, e)
	row := e.inputModeLookupTable.Candidates[config.UsIM-1].Value().(ibus.Text)
	if row.Text != config.UsIMRemoveLabel {
		t.Fatalf("exclusion label = %q, want %q", row.Text, config.UsIMRemoveLabel)
	}
	selectMode(t, e, '7')
	if got := e.getInputMode(); got != config.PreeditIM {
		t.Fatalf("mode after removing the exclusion = %d, want %d", got, config.PreeditIM)
	}
	if got := int(atomic.LoadInt32(&sessionInputMode)); got != 0 {
		t.Fatalf("session mode = %d, want cleared", got)
	}

	const wmClass = "org.gnome.Ptyxis"
	e.wmClasses = wmClass
	e.engineName = "input-mode-exclusion-test"
	e.config.InputModeMapping[wmClass] = config.UsIM
	defer os.Remove(config.GetConfigPath(e.engineName))
	openModeTable(t, e)
	labeled := e.inputModeLookupTable.Candidates[config.UsIM-1].Value().(ibus.Text)
	if labeled.Text != config.UsIMRemoveLabel+" ("+wmClass+")" {
		t.Fatalf("per-app exclusion label = %q", labeled.Text)
	}
	selectMode(t, e, '7')
	if _, still := e.config.InputModeMapping[wmClass]; still {
		t.Fatalf("excluded app stayed in the map: %#v", e.config.InputModeMapping)
	}
	if got := e.getInputMode(); got != config.PreeditIM {
		t.Fatalf("mode for an unlisted app = %d, want %d", got, config.PreeditIM)
	}
}

func openModeTable(t *testing.T, e *IBusBambooEngine) {
	t.Helper()
	handled, err := e.ProcessKeyEvent('~', '~', IBusShiftMask)
	if err != nil || !handled || !e.isInputModeLTOpened {
		t.Fatalf("Shift+~ should open the mode table: handled=%v opened=%v err=%v", handled, e.isInputModeLTOpened, err)
	}
}

func selectMode(t *testing.T, e *IBusBambooEngine, key rune) {
	t.Helper()
	handled, err := e.ProcessKeyEvent(uint32(key), uint32(key), 0)
	if err != nil || !handled || e.isInputModeLTOpened {
		t.Fatalf("selecting mode %c failed: handled=%v opened=%v err=%v", key, handled, e.isInputModeLTOpened, err)
	}
}

func TestInputModeSwitchKeepsPerAppMapping(t *testing.T) {
	atomic.StoreInt32(&sessionInputMode, int32(config.SurroundingTextIM))
	defer atomic.StoreInt32(&sessionInputMode, 0)

	const wmClass = "org.gnome.TextEditor"
	e := newTestEngine(config.PreeditIM)
	e.engineName = "input-mode-fallback-test"
	e.wmClasses = wmClass
	defer os.Remove(config.GetConfigPath(e.engineName))

	handled, err := e.ProcessKeyEvent('~', '~', IBusShiftMask)
	if err != nil || !handled || !e.isInputModeLTOpened {
		t.Fatalf("Shift+~ should open the mode table when a wm class is known: handled=%v opened=%v err=%v", handled, e.isInputModeLTOpened, err)
	}
	handled, err = e.ProcessKeyEvent('3', '3', 0)
	if err != nil || !handled {
		t.Fatalf("selecting a per-application mode failed: handled=%v err=%v", handled, err)
	}
	if e.config.InputModeMapping[wmClass] != config.BackspaceForwardingIM {
		t.Fatalf("mapped mode = %d, want %d", e.config.InputModeMapping[wmClass], config.BackspaceForwardingIM)
	}
	if got := int(atomic.LoadInt32(&sessionInputMode)); got != config.SurroundingTextIM {
		t.Fatalf("session mode = %d, want %d", got, config.SurroundingTextIM)
	}
	if got := e.getInputMode(); got != config.BackspaceForwardingIM {
		t.Fatalf("mode for a known class = %d, want %d", got, config.BackspaceForwardingIM)
	}

	e.wmClasses = "other.app"
	if got := e.getInputMode(); got != config.PreeditIM {
		t.Fatalf("mode for an unmapped class = %d, want default %d", got, config.PreeditIM)
	}
}

func newTestEngine(inputMode int) *IBusBambooEngine {
	fe := NewFakeEngine()
	cfg := config.DefaultCfg()
	cfg.DefaultInputMode = inputMode
	inputMethod := bamboo.ParseInputMethod(cfg.InputMethodDefinitions, cfg.InputMethod)
	return NewIbusBambooEngine("test", &cfg, fe, bamboo.NewEngine(inputMethod, cfg.Flags))
}

func assertEngine(t testing.TB, tc testCase, assertFn func(testing.TB, *fakeEngine, IEngine)) {
	fe := NewFakeEngine()
	engineName := "test"
	var cfg = config.DefaultCfg()
	cfg.DefaultInputMode = tc.inputMode
	inputMethod := bamboo.ParseInputMethod(cfg.InputMethodDefinitions, cfg.InputMethod)
	if tc.mTable != nil {
		cfg.IBflags |= config.IBmacroEnabled
	}
	e := NewIbusBambooEngine(engineName, &cfg, fe, bamboo.NewEngine(inputMethod, cfg.Flags))
	if tc.mTable != nil {
		e.macroTable = &MacroTable{
			mTable: tc.mTable,
		}
	}
	assertFn(t, fe, e)
}
