package ibus

import (
	"github.com/godbus/dbus"
)

type LookupTable struct {
	Name          string
	Attachments   map[string]dbus.Variant
	PageSize      uint32
	CursorPos     uint32
	CursorVisible bool
	Round         bool
	Orientation   int32
	Candidates    []dbus.Variant
	Labels        []dbus.Variant
}

func NewLookupTable() *LookupTable {
	lt := &LookupTable{}
	lt.Name = "IBusLookupTable"
	lt.PageSize = 5
	lt.CursorPos = 0
	lt.CursorVisible = true
	lt.Round = false
	lt.Orientation = ORIENTATION_SYSTEM

	return lt
}

func (lt *LookupTable) AppendCandidate(text string) {
	t := NewText(text)
	lt.Candidates = append(lt.Candidates, dbus.MakeVariant(*t))
}

func (lt *LookupTable) AppendLabel(label string) {
	l := NewText(label)
	lt.Labels = append(lt.Labels, dbus.MakeVariant(*l))
}
