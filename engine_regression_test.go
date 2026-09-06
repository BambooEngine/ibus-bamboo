package main

import (
	"ibus-bamboo/config"
	"sync"
	"testing"
	"time"

	"github.com/BambooEngine/bamboo-core"
	ibus "github.com/BambooEngine/goibus"
	"github.com/godbus/dbus/v5"
)

func newRegressionEngine(inputMode int) (*IBusBambooEngine, *fakeEngine) {
	fe := NewFakeEngine()
	cfg := config.DefaultCfg()
	cfg.DefaultInputMode = inputMode
	im := bamboo.ParseInputMethod(cfg.InputMethodDefinitions, cfg.InputMethod)
	return NewIbusBambooEngine("test", &cfg, fe, bamboo.NewEngine(im, cfg.Flags)), fe
}

func TestResetClearsBackspaceBuffer(t *testing.T) {
	e, _ := newRegressionEngine(config.SurroundingTextIM)
	e.ProcessKeyEvent('d', 'd', 0)
	e.ProcessKeyEvent('u', 'u', 0)
	if e.getRawKeyLen() == 0 {
		t.Fatalf("setup failed: buffer should be non-empty before Reset")
	}
	e.Reset()
	if got := e.getRawKeyLen(); got != 0 {
		t.Errorf("after Reset rawKeyLen = %d, want 0", got)
	}
}

func TestResetClearsPreeditBuffer(t *testing.T) {
	e, fe := newRegressionEngine(config.PreeditIM)
	e.ProcessKeyEvent('d', 'd', 0)
	e.ProcessKeyEvent('u', 'u', 0)
	e.Reset()
	if got := e.getRawKeyLen(); got != 0 {
		t.Errorf("after Reset rawKeyLen = %d, want 0", got)
	}
	if !fe.getHidePreeditText() {
		t.Errorf("after Reset preedit should be hidden")
	}
}

func TestFocusOutCommitsToLosingWindow(t *testing.T) {
	e, fe := newRegressionEngine(config.PreeditIM)
	e.ProcessKeyEvent('d', 'd', 0)
	e.ProcessKeyEvent('u', 'u', 0)
	pending := fe.getPreeditText()
	if pending == "" {
		t.Fatalf("setup failed: expected pending preedit")
	}
	e.FocusOut()
	if fe.getCommitText() != pending {
		t.Errorf("after FocusOut commit = %q, want pending preedit %q", fe.getCommitText(), pending)
	}
	if got := e.getRawKeyLen(); got != 0 {
		t.Errorf("after FocusOut rawKeyLen = %d, want 0", got)
	}
}

func TestFocusOutClearsBackspaceBuffer(t *testing.T) {
	e, _ := newRegressionEngine(config.SurroundingTextIM)
	e.ProcessKeyEvent('d', 'd', 0)
	e.FocusOut()
	if got := e.getRawKeyLen(); got != 0 {
		t.Errorf("after FocusOut rawKeyLen = %d, want 0", got)
	}
}

func TestMovementHomeResetsBuffer(t *testing.T) {
	e, _ := newRegressionEngine(config.SurroundingTextIM)
	e.ProcessKeyEvent('d', 'd', 0)
	if e.getRawKeyLen() == 0 {
		t.Fatalf("setup failed: buffer should be non-empty")
	}
	if ret, _ := e.ProcessKeyEvent(IBusHome, 0, 0); ret {
		t.Errorf("Home should be forwarded (return false), got true")
	}
	if got := e.getRawKeyLen(); got != 0 {
		t.Errorf("after Home rawKeyLen = %d, want 0", got)
	}
}

func TestExtractSurroundingString(t *testing.T) {
	if s, ok := extractSurroundingString("abc"); !ok || s != "abc" {
		t.Errorf("plain string: got %q,%v", s, ok)
	}
	if s, ok := extractSurroundingString(dbus.MakeVariant("abc").Value()); !ok || s != "abc" {
		t.Errorf("variant string: got %q,%v", s, ok)
	}
	if s, ok := extractSurroundingString([]interface{}{"IBusText", 0, "abc"}); !ok || s != "abc" {
		t.Errorf("tuple shape: got %q,%v", s, ok)
	}
	if s, ok := extractSurroundingString(*ibus.NewText("abc")); !ok || s != "abc" {
		t.Errorf("goibus.Text: got %q,%v", s, ok)
	}
	if _, ok := extractSurroundingString(42); ok {
		t.Errorf("garbage shape should not decode")
	}
	if _, ok := extractSurroundingString(nil); ok {
		t.Errorf("nil should not decode")
	}
}

func TestSetSurroundingTextRebuildsBuffer(t *testing.T) {
	e, _ := newRegressionEngine(config.SurroundingTextIM)
	e.surroundingTextReady = true
	e.SetSurroundingText(dbus.MakeVariant(*ibus.NewText("hello")), 5, 5)
	if got := e.getRawKeyLen(); got != 5 {
		t.Errorf("after SetSurroundingText rawKeyLen = %d, want 5", got)
	}
	if e.surroundingTextReady {
		t.Errorf("ready flag should be consumed")
	}
	// A selection (anchor != cursor) must be ignored.
	e.surroundingTextReady = true
	e.SetSurroundingText(dbus.MakeVariant(*ibus.NewText("hello world")), 5, 2)
	if got := e.getRawKeyLen(); got != 5 {
		t.Errorf("selection must be ignored, rawKeyLen = %d, want 5", got)
	}
}

func TestPerEngineQueueIsolation(t *testing.T) {
	e1, fe1 := newRegressionEngine(config.SurroundingTextIM)
	e2, fe2 := newRegressionEngine(config.SurroundingTextIM)
	e1.shouldEnqueuKeyStrokes = true
	e2.shouldEnqueuKeyStrokes = true
	if ok, _ := e1.ProcessKeyEvent('d', 'd', 0); !ok {
		t.Fatalf("key not accepted")
	}
	if ok, _ := e1.ProcessKeyEvent('u', 'u', 0); !ok {
		t.Fatalf("key not accepted")
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if fe1.getCommitText() == "du" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if fe1.getCommitText() != "du" {
		t.Errorf("engine1 commit = %q, want %q", fe1.getCommitText(), "du")
	}
	if fe2.getCommitText() != "" {
		t.Errorf("engine1 keystrokes leaked onto engine2: %q", fe2.getCommitText())
	}
	if got := e2.getRawKeyLen(); got != 0 {
		t.Errorf("engine2 buffer len = %d, want 0", got)
	}
}

func TestConcurrentEngineHammer(t *testing.T) {
	e, _ := newRegressionEngine(config.SurroundingTextIM)
	e.shouldEnqueuKeyStrokes = true
	var wg sync.WaitGroup
	keys := []rune(" engine ")
	for g := 0; g < 3; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i, k := range keys {
				e.ProcessKeyEvent(uint32(k), uint32(k), 0)
				if i%3 == 0 {
					e.Reset()
				}
			}
			e.ProcessKeyEvent(IBusBackSpace, 0, 0)
			e.FocusOut()
			e.SetCapabilities(0)
			e.PageUp()
			e.CandidateClicked(0, 1, 0)
		}(g)
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(20 * time.Second):
		t.Fatalf("deadlock under concurrent engine use")
	}
}
