package ibus

import (
	"github.com/godbus/dbus"
)

type EngineDesc struct {
	Name          string                  `xml:"-"`
	Attachments   map[string]dbus.Variant `xml:"-"`
	EngineName    string                  `xml:"name"`
	LongName      string                  `xml:"longname"`
	Description   string                  `xml:"description"`
	Language      string                  `xml:"language"`
	License       string                  `xml:"license"`
	Author        string                  `xml:"author"`
	Icon          string                  `xml:"icon"`
	Layout        string                  `xml:"layout"`
	Rank          uint32                  `xml:"rank"`
	Hotkeys       string                  `xml:"hotkeys,omitempty"`
	Symbol        string                  `xml:"symbol,omitempty"`
	Setup         string                  `xml:"setup,omitempty"`
	LayoutVariant string                  `xml:"layout-variant,omitempty"`
	LayoutOption  string                  `xml:"layout-option,omitempty"`
	Version       string                  `xml:"version,omitempty"`
	Textdomain    string                  `xml:"textdomain,omitempty"`
}

func TinyEngineDesc(name string, longname string, desc string, lang string, license string, author string, icon string, layout string) *EngineDesc {
	ed := &EngineDesc{}

	ed.Name = "IBusEngineDesc"
	ed.EngineName = name
	ed.LongName = longname
	ed.Description = desc
	ed.Language = lang
	ed.License = license
	ed.Author = author
	ed.Icon = icon
	ed.Layout = layout

	return ed
}

func SmallEngineDesc(name string, longname string, desc string, lang string, license string, author string, icon string, layout string,
	setup string, version string) *EngineDesc {

	ed := TinyEngineDesc(name, longname, desc, lang, license, author, icon, layout)
	ed.Setup = setup
	ed.Version = version

	return ed
}

func FullEngineDesc(name string, longname string, desc string, lang string, license string, author string, icon string, layout string,
	rank uint32, hotkeys string, symbol string, setup string, layoutVariant string, layoutOption string, version string, textdomain string) *EngineDesc {

	ed := TinyEngineDesc(name, longname, desc, lang, license, author, icon, layout)
	ed.Rank = rank
	ed.Hotkeys = hotkeys
	ed.Symbol = symbol
	ed.Setup = setup
	ed.LayoutVariant = layoutVariant
	ed.LayoutOption = layoutOption
	ed.Version = version
	ed.Textdomain = textdomain

	return ed
}
