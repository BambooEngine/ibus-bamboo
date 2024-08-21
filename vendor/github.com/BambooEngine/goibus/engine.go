package goibus

import (
	"github.com/godbus/dbus/v5"
)

type Engine struct {
	conn       *dbus.Conn
	objectPath dbus.ObjectPath
}

func BaseEngine(conn *dbus.Conn, objectPath dbus.ObjectPath) Engine {
	return Engine{conn, objectPath}
}

func PublishEngine(conn *dbus.Conn, objectPath dbus.ObjectPath, userEngine interface{}) {
	conn.Export(userEngine, objectPath, IBUS_IFACE_ENGINE)
	conn.Export(userEngine, objectPath, IBUS_IFACE_SERVICE)
	conn.Export(userEngine, objectPath, BUS_PROPERTIES_NAME)
}

func (e *Engine) emitSignal(name string, values ...interface{}) {
	methodName := IBUS_IFACE_ENGINE + "." + name
	e.conn.Emit(e.objectPath, methodName, values...)
}

func (e *Engine) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	items := make(map[string]dbus.Variant)
	return items, nil
}

//@method(in_signature="uuu", out_signature="b")
func (e *Engine) ProcessKeyEvent(keyval uint32, keycode uint32, state uint32) (bool, *dbus.Error) {
	return false, nil
}

//@method(in_signature="iiii")
func (e *Engine) SetCursorLocation(x int32, y int32, w int32, h int32) *dbus.Error {
	return nil
}

//@method(in_signature="vuu")
func (e *Engine) SetSurroundingText(text dbus.Variant, cursor_index uint32, anchor_pos uint32) *dbus.Error {
	return nil
}

//@method(in_signature="u")
func (e *Engine) SetCapabilities(cap uint32) *dbus.Error {
	return nil
}

//@method()
func (e *Engine) FocusIn() *dbus.Error {
	return nil
}

//@method()
func (e *Engine) FocusOut() *dbus.Error {
	return nil
}

//@method()
func (e *Engine) Reset() *dbus.Error {
	return nil
}

//@method()
func (e *Engine) PageUp() *dbus.Error {
	return nil
}

//@method()
func (e *Engine) PageDown() *dbus.Error {
	return nil
}

//@method()
func (e *Engine) CursorUp() *dbus.Error {
	return nil
}

//@method()
func (e *Engine) CursorDown() *dbus.Error {
	return nil
}

//@method(in_signature="uuu")
func (e *Engine) CandidateClicked(index uint32, button uint32, state uint32) *dbus.Error {
	return nil
}

//@method()
func (e *Engine) Enable() *dbus.Error {
	return nil
}

//@method()
func (e *Engine) Disable() *dbus.Error {
	return nil
}

//@method(in_signature="su")
func (e *Engine) PropertyActivate(prop_name string, prop_state uint32) *dbus.Error {
	return nil
}

//@method(in_signature="s")
func (e *Engine) PropertyShow(prop_name string) *dbus.Error {
	return nil
}

//@method(in_signature="s")
func (e *Engine) PropertyHide(prop_name string) *dbus.Error {
	return nil
}

//@method()
func (e *Engine) Destroy() *dbus.Error {
	e.conn.Export(nil, e.objectPath, IBUS_IFACE_ENGINE)
	e.conn.Export(nil, e.objectPath, IBUS_IFACE_SERVICE)
	e.conn.Export(nil, e.objectPath, BUS_PROPERTIES_NAME)
	return nil
}

//@signal(signature="v")
func (e *Engine) CommitText(text *Text) {
	e.emitSignal("CommitText", dbus.MakeVariant(*text))
}

//@signal(signature="uuu")
func (e *Engine) ForwardKeyEvent(keyval uint32, keycode uint32, state uint32) {
	e.emitSignal("ForwardKeyEvent", keyval, keycode, state)
}

//@signal(signature="vubu")
func (e *Engine) UpdatePreeditText(text *Text, cursor_pos uint32, visible bool) {
	e.emitSignal("UpdatePreeditText", dbus.MakeVariant(*text), cursor_pos, visible, IBUS_ENGINE_PREEDIT_CLEAR)
}
func (e *Engine) UpdatePreeditTextWithMode(text *Text, cursor_pos uint32, visible bool, mode uint32) {
	e.emitSignal("UpdatePreeditText", dbus.MakeVariant(*text), cursor_pos, visible, mode)
}

//@signal()
func (e *Engine) ShowPreeditText() {
	e.emitSignal("ShowPreeditText")
}

//@signal()
func (e *Engine) HidePreeditText() {
	e.emitSignal("HidePreeditText")
}

//@signal(signature="vb")
func (e *Engine) UpdateAuxiliaryText(text *Text, visible bool) {
	e.emitSignal("UpdateAuxiliaryText", dbus.MakeVariant(*text), visible)
}

//@signal()
func (e *Engine) ShowAuxiliaryText() {
	e.emitSignal("ShowAuxiliaryText")
}

//@signal()
func (e *Engine) HideAuxiliaryText() {
	e.emitSignal("HideAuxiliaryText")
}

//@signal(signature="vb")
func (e *Engine) UpdateLookupTable(lookup_table *LookupTable, visible bool) {
	e.emitSignal("UpdateLookupTable", dbus.MakeVariant(*lookup_table), visible)
}

//@signal()
func (e *Engine) ShowLookupTable() {
	e.emitSignal("ShowLookupTable")
}

//@signal()
func (e *Engine) HideLookupTable() {
	e.emitSignal("HideLookupTable")
}

//@signal()
func (e *Engine) PageUpLookupTable() {
	e.emitSignal("PageUpLookupTable")
}

//@signal()
func (e *Engine) PageDownLookupTable() {
	e.emitSignal("PageDownLookupTable")
}

//@signal()
func (e *Engine) CursorUpLookupTable() {
	e.emitSignal("CursorUpLookupTable")
}

//@signal()
func (e *Engine) CursorDownLookupTable() {
	e.emitSignal("CursorDownLookupTable")
}

//@signal(signature="v")
func (e *Engine) RegisterProperties(props *PropList) {
	e.emitSignal("RegisterProperties", dbus.MakeVariant(*props))
}

//@signal(signature="v")
func (e *Engine) UpdateProperty(prop *Property) {
	e.emitSignal("UpdateProperty", dbus.MakeVariant(*prop))
}

//@signal(signature="iu")
func (e *Engine) DeleteSurroundingText(offset_from_cursor int32, nchars uint32) {
	e.emitSignal("DeleteSurroundingText", offset_from_cursor, nchars)
}

//@signal()
func (e *Engine) RequireSurroundingText() {
	e.emitSignal("RequireSurroundingText")
}
