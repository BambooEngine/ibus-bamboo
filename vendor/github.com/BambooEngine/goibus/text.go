package goibus

import (
	"github.com/godbus/dbus/v5"
)

type Attribute struct {
	Name        string
	Attachments map[string]dbus.Variant
	Type        uint32
	Value       uint32
	StartIndex  uint32
	EndIndex    uint32
}

type AttrList struct {
	Name        string
	Attachments map[string]dbus.Variant
	Attributes  []dbus.Variant
}

type Text struct {
	Name        string
	Attachments map[string]dbus.Variant
	Text        string
	AttrList    dbus.Variant
}

func NewAttribute(attrType, attrValue, startIndex uint32, endIndex uint32) *Attribute {
	var attr = Attribute{
		Name:       "IBusAttribute",
		Type:       attrType,
		Value:      attrValue,
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}
	return &attr
}

func (t *Text) AppendAttr(attrType, attrValue, startIndex uint32, endIndex uint32) {
	attrList := &AttrList{}
	attrList.Name = "IBusAttrList"
	attrList.Attributes = append(attrList.Attributes, dbus.MakeVariant(*NewAttribute(attrType, attrValue, startIndex, endIndex)))
	t.AttrList = dbus.MakeVariant(*attrList)
}

func NewText(text string) *Text {
	attrList := AttrList{}
	attrList.Name = "IBusAttrList"

	t := Text{}
	t.Name = "IBusText"
	t.Text = text
	t.AttrList = dbus.MakeVariant(attrList)

	return &t
}
