package main

import (
	"sync"

	ibus "github.com/BambooEngine/goibus"
	"github.com/godbus/dbus/v5"
)

type IEngine interface {
	GetAll(iface string) (map[string]dbus.Variant, *dbus.Error)
	ProcessKeyEvent(keyval uint32, keycode uint32, state uint32) (bool, *dbus.Error)
	SetCursorLocation(x int32, y int32, w int32, h int32) *dbus.Error
	SetSurroundingText(text dbus.Variant, cursor_index uint32, anchor_pos uint32) *dbus.Error
	SetCapabilities(cap uint32) *dbus.Error
	FocusIn() *dbus.Error
	FocusOut() *dbus.Error
	Reset() *dbus.Error
	PageUp() *dbus.Error
	PageDown() *dbus.Error
	CursorUp() *dbus.Error
	CursorDown() *dbus.Error
	CandidateClicked(index uint32, button uint32, state uint32) *dbus.Error
	Enable() *dbus.Error
	Disable() *dbus.Error
	PropertyActivate(prop_name string, prop_state uint32) *dbus.Error
	PropertyShow(prop_name string) *dbus.Error
	PropertyHide(prop_name string) *dbus.Error
	Destroy() *dbus.Error
	CommitText(text *ibus.Text)
	ForwardKeyEvent(keyval uint32, keycode uint32, state uint32)
	UpdatePreeditText(text *ibus.Text, cursor_pos uint32, visible bool)
	UpdatePreeditTextWithMode(text *ibus.Text, cursor_pos uint32, visible bool, mode uint32)
	ShowPreeditText()
	HidePreeditText()
	UpdateAuxiliaryText(text *ibus.Text, visible bool)
	ShowAuxiliaryText()
	HideAuxiliaryText()
	UpdateLookupTable(lookup_table *ibus.LookupTable, visible bool)
	ShowLookupTable()
	HideLookupTable()
	PageUpLookupTable()
	PageDownLookupTable()
	CursorUpLookupTable()
	CursorDownLookupTable()
	RegisterProperties(props *ibus.PropList)
	UpdateProperty(prop *ibus.Property)
	DeleteSurroundingText(offset_from_cursor int32, nchars uint32)
	RequireSurroundingText()
}

// fakeEngine is used by tests; its assertions may run on a different
// goroutine than the engine's key queue consumer, so state is guarded.
type fakeEngine struct {
	mu                  sync.Mutex
	commitText          string
	preeditText         string
	committed           bool
	isHidePreeditText   bool
	isHideAuxiliaryText bool
	isHideLookupTable   bool
	isReset             bool
	forwardKeyEvent     [3]uint32
}

func (e *fakeEngine) getCommitText() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.commitText
}

func (e *fakeEngine) getPreeditText() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.preeditText
}

func (e *fakeEngine) getHidePreeditText() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.isHidePreeditText
}

func NewFakeEngine() *fakeEngine {
	return &fakeEngine{}
}

func (e *fakeEngine) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	items := make(map[string]dbus.Variant)
	return items, nil
}

func (e *fakeEngine) ProcessKeyEvent(keyval uint32, keycode uint32, state uint32) (bool, *dbus.Error) {
	return false, nil
}

func (e *fakeEngine) SetCursorLocation(x int32, y int32, w int32, h int32) *dbus.Error {
	return nil
}

func (e *fakeEngine) SetSurroundingText(text dbus.Variant, cursor_index uint32, anchor_pos uint32) *dbus.Error {
	return nil
}

func (e *fakeEngine) SetCapabilities(cap uint32) *dbus.Error {
	return nil
}

func (e *fakeEngine) FocusIn() *dbus.Error {
	return nil
}

func (e *fakeEngine) FocusOut() *dbus.Error {
	return nil
}

func (e *fakeEngine) Reset() *dbus.Error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isReset = true
	return nil
}

// @method()
func (e *fakeEngine) PageUp() *dbus.Error {
	return nil
}

// @method()
func (e *fakeEngine) PageDown() *dbus.Error {
	return nil
}

// @method()
func (e *fakeEngine) CursorUp() *dbus.Error {
	return nil
}

// @method()
func (e *fakeEngine) CursorDown() *dbus.Error {
	return nil
}

// @method(in_signature="uuu")
func (e *fakeEngine) CandidateClicked(index uint32, button uint32, state uint32) *dbus.Error {
	return nil
}

// @method()
func (e *fakeEngine) Enable() *dbus.Error {
	return nil
}

// @method()
func (e *fakeEngine) Disable() *dbus.Error {
	return nil
}

// @method(in_signature="su")
func (e *fakeEngine) PropertyActivate(prop_name string, prop_state uint32) *dbus.Error {
	return nil
}

// @method(in_signature="s")
func (e *fakeEngine) PropertyShow(prop_name string) *dbus.Error {
	return nil
}

// @method(in_signature="s")
func (e *fakeEngine) PropertyHide(prop_name string) *dbus.Error {
	return nil
}

// @method()
func (e *fakeEngine) Destroy() *dbus.Error {
	return nil
}

// @signal(signature="v")
func (e *fakeEngine) CommitText(text *ibus.Text) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.commitText += text.Text
}

// @signal(signature="uuu")
func (e *fakeEngine) ForwardKeyEvent(keyval uint32, keycode uint32, state uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.forwardKeyEvent = [3]uint32{keyval, keycode, state}
}

// @signal(signature="vubu")
func (e *fakeEngine) UpdatePreeditText(text *ibus.Text, cursor_pos uint32, visible bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.preeditText = text.Text
}
func (e *fakeEngine) UpdatePreeditTextWithMode(text *ibus.Text, cursor_pos uint32, visible bool, mode uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.preeditText = text.Text
}

// @signal()
func (e *fakeEngine) ShowPreeditText() {
}

// @signal()
func (e *fakeEngine) HidePreeditText() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.preeditText = ""
	e.isHidePreeditText = true
}

// @signal(signature="vb")
func (e *fakeEngine) UpdateAuxiliaryText(text *ibus.Text, visible bool) {
}

// @signal()
func (e *fakeEngine) ShowAuxiliaryText() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isHideAuxiliaryText = false
}

// @signal()
func (e *fakeEngine) HideAuxiliaryText() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isHideAuxiliaryText = true
}

// @signal(signature="vb")
func (e *fakeEngine) UpdateLookupTable(lookup_table *ibus.LookupTable, visible bool) {
}

// @signal()
func (e *fakeEngine) ShowLookupTable() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isHideLookupTable = false
}

// @signal()
func (e *fakeEngine) HideLookupTable() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isHideLookupTable = true
}

// @signal()
func (e *fakeEngine) PageUpLookupTable() {
}

// @signal()
func (e *fakeEngine) PageDownLookupTable() {
}

// @signal()
func (e *fakeEngine) CursorUpLookupTable() {
}

// @signal()
func (e *fakeEngine) CursorDownLookupTable() {
}

// @signal(signature="v")
func (e *fakeEngine) RegisterProperties(props *ibus.PropList) {
}

// @signal(signature="v")
func (e *fakeEngine) UpdateProperty(prop *ibus.Property) {
}

// @signal(signature="iu")
func (e *fakeEngine) DeleteSurroundingText(offset_from_cursor int32, nchars uint32) {
	s := []rune(e.commitText)
	var txt string
	for _, ch := range s[:len(s)-int(nchars)] {
		txt += string(ch)
	}
	e.commitText = txt
}

// @signal()
func (e *fakeEngine) RequireSurroundingText() {
}
