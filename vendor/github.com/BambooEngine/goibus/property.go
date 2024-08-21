package goibus

import (
	"github.com/godbus/dbus/v5"
)

type PropList struct {
	Name         string
	Attachments  map[string]dbus.Variant
	PropertyList []dbus.Variant
}

type Property struct {
	Name        string
	Attachments map[string]dbus.Variant
	Key         string
	Type        uint32
	Label       dbus.Variant
	Icon        string
	Tooltip     dbus.Variant
	Sensitive   bool
	Visible     bool
	State       uint32
	SubProps    dbus.Variant
	Symbol      dbus.Variant
}

func NewProperty(key string, ptype uint32, label string, icon string, tooltip string, sensitive bool, visible bool, state uint32) *Property {
	p := &Property{}
	p.Name = "IBusProperty"
	p.Key = key
	p.Type = ptype
	p.Label = dbus.MakeVariant(*NewText(label))
	p.Icon = icon
	p.Tooltip = dbus.MakeVariant(*NewText(tooltip))
	p.Sensitive = sensitive
	p.Visible = visible
	p.State = state
	p.SubProps = dbus.MakeVariant(*NewPropList())
	p.Symbol = dbus.MakeVariant(*NewText(""))

	return p
}

func NewPropertyWithChild(key string, ptype uint32, label string, icon string, tooltip string, sensitive bool, visible bool, state uint32, child PropList) *Property {
	p := NewProperty(key, ptype, label, icon, tooltip, sensitive, visible, state)
	p.SubProps = dbus.MakeVariant(child)

	return p
}

func NewPropList(propList ...*Property) *PropList {
	pl := &PropList{}
	pl.Name = "IBusPropList"

	for _, p := range propList {
		pl.PropertyList = append(pl.PropertyList, dbus.MakeVariant(*p))
	}

	return pl
}
