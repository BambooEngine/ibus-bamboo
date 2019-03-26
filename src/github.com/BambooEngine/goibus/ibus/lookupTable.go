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

func (lt *LookupTable) SetCursorPos(pos uint32) bool {
	if pos >= uint32(len(lt.Candidates)) || pos < 0 {
		return false
	}
	lt.CursorPos = pos
	return true
}

func (lt *LookupTable) GetCursorPos() uint32 {
	return lt.CursorPos
}

func (lt *LookupTable) GetCursorPosInCurrentPage() uint32 {
	return lt.CursorPos % lt.PageSize
}

func (lt *LookupTable) SetCursorPosInCurrentPage(pos uint32) bool {
	if pos < 0 || pos >= lt.PageSize {
		return false
	}
	pos += lt.GetCursorPosInCurrentPage()
	if pos >= uint32(len(lt.Candidates)) {
		return false
	}
	lt.CursorPos = pos
	return true
}

func (lt *LookupTable) CursorUp() bool {
	if lt.CursorPos == 0 {
		if lt.Round {
			lt.CursorPos = uint32(len(lt.Candidates)) - 1
			return true
		} else {
			return false
		}
	}
	lt.CursorPos -= 1
	return true
}

func (lt *LookupTable) CursorDown() bool {
	if lt.CursorPos == uint32(len(lt.Candidates)) {
		if lt.Round {
			lt.CursorPos = 0
			return true
		} else {
			return false
		}
	}
	lt.CursorPos += 1
	return true
}

func (lt *LookupTable) PageUp() bool {
	if lt.CursorPos < lt.PageSize {
		if lt.Round {
			var nrCandidates = uint32(len(lt.Candidates))
			var maxPage = uint32(nrCandidates / lt.PageSize)
			lt.CursorPos += maxPage * lt.PageSize
			if lt.CursorPos > nrCandidates-1 {
				lt.CursorPos = nrCandidates - 1
				return true
			} else {
				return false
			}
		}
	}
	lt.CursorPos -= lt.PageSize
	return true
}

func (lt *LookupTable) PageDown() bool {
	var currentPage = lt.CursorPos / lt.PageSize
	var nrCandidates = uint32(len(lt.Candidates))
	var maxPage = nrCandidates / lt.PageSize
	if currentPage >= maxPage {
		if lt.Round {
			lt.CursorPos %= lt.PageSize
			return true
		} else {
			return false
		}
	}
	var pos = lt.CursorPos + lt.PageSize
	if pos >= nrCandidates {
		pos = nrCandidates - 1
	}
	lt.CursorPos = pos
	return true
}

func (lt *LookupTable) Clean() {
	lt.Candidates = nil
	lt.Labels = nil
	lt.CursorPos = 0
}
