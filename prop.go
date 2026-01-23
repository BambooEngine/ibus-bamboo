/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"ibus-bamboo/config"

	"github.com/BambooEngine/bamboo-core"
	ibus "github.com/BambooEngine/goibus"
	"github.com/godbus/dbus/v5"
)

const (
	PropKeyAbout                        = "about"
	PropKeyStdToneStyle                 = "std_tone_style"
	PropKeyFreeToneMarking              = "tone_free_marking"
	PropKeyEnableSpellCheck             = "enable_spell_check"
	PropKeySpellCheckByRules            = "spell_check_by_rules"
	PropKeySpellCheckByDicts            = "spell_check_by_dicts"
	PropKeyPreeditInvisibility          = "preedit_invisibility"
	PropKeyVnCharsetConvert             = "charset_convert_page"
	PropKeyMacroEnabled                 = "macro_enabled"
	PropKeyMacroTable                   = "open_macro_table"
	PropKeyEmojiEnabled                 = "emoji_enabled"
	PropKeyConfiguration                = "configuration"
	PropKeyPreeditElimination           = "preedit_elimination"
	PropKeyInputModeLookupTable         = "input_mode_lookup_table"
	PropKeyInputModeLookupTableShortcut = "input_mode_lookup_table_shortcut"
	PropKeyAutoCapitalizeMacro          = "auto_capitalize_macro"
	PropKeyIMQuickSwitchEnabled         = "im_quick_switch"
	PropKeyRestoreKeyStrokes            = "restore_key_strokes"
)

var IBusSeparator = &ibus.Property{
	Name:      "IBusProperty",
	Key:       "-",
	Type:      ibus.PROP_TYPE_SEPARATOR,
	Label:     dbus.MakeVariant(ibus.NewText("")),
	Tooltip:   dbus.MakeVariant(ibus.NewText("")),
	Sensitive: true,
	Visible:   true,
	Symbol:    dbus.MakeVariant(ibus.NewText("")),
	SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
}

func GetPropListByConfig(c *config.Config) *ibus.PropList {
	var aboutText = "IBus " + EngineName + " " + Version
	if !*embedded {
		aboutText += " (Debug)"
	}
	if c.DefaultInputMode == config.UsIM {
		return ibus.NewPropList(
			&ibus.Property{
				Name:      "IBusProperty",
				Key:       PropKeyAbout,
				Type:      ibus.PROP_TYPE_NORMAL,
				Label:     dbus.MakeVariant(ibus.NewText(aboutText)),
				Tooltip:   dbus.MakeVariant(ibus.NewText("Mở trang chủ")),
				Sensitive: true,
				Visible:   true,
				Icon:      "gtk-home",
				Symbol:    dbus.MakeVariant(ibus.NewText("")),
				SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
			},
			&ibus.Property{
				Name:      "IBusProperty",
				Key:       PropKeyInputModeLookupTableShortcut,
				Type:      ibus.PROP_TYPE_NORMAL,
				Label:     dbus.MakeVariant(ibus.NewText("Keyboard Shortcuts")),
				Tooltip:   dbus.MakeVariant(ibus.NewText("Keyboard Shortcuts")),
				Sensitive: true,
				Visible:   true,
				Icon:      "appointment",
				Symbol:    dbus.MakeVariant(ibus.NewText("")),
				SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
			},
		)
	}
	return ibus.NewPropList(
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAbout,
			Type:      ibus.PROP_TYPE_NORMAL,
			Label:     dbus.MakeVariant(ibus.NewText(aboutText)),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Mở trang chủ")),
			Sensitive: true,
			Visible:   true,
			Icon:      "gtk-home",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		IBusSeparator,
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_MENU,
			Label:     dbus.MakeVariant(ibus.NewText("Bảng mã")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Bảng mã")),
			Sensitive: true,
			Visible:   true,
			Icon:      "fonts",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(GetCharsetPropListByConfig(c)),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_MENU,
			Label:     dbus.MakeVariant(ibus.NewText("Kiểu gõ")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Kiểu gõ")),
			Sensitive: true,
			Visible:   true,
			Icon:      "preferences-desktop",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(GetIMPropListByConfig(c)),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_MENU,
			Label:     dbus.MakeVariant(ibus.NewText("Gõ tắt")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Gõ tắt")),
			Sensitive: true,
			Visible:   true,
			Icon:      "document-send",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(GetMacroPropListByConfig(c)),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_MENU,
			Label:     dbus.MakeVariant(ibus.NewText("Kiểm tra chính tả")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Kiểm tra chính tả")),
			Sensitive: true,
			Visible:   true,
			Icon:      "tools-check-spelling",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(GetSpellCheckingPropListByConfig(c)),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_MENU,
			Label:     dbus.MakeVariant(ibus.NewText("Cấu hình khác")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Cấu hình khác")),
			Sensitive: true,
			Visible:   true,
			Icon:      "preferences-other",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(GetOptionsPropListByConfig(c)),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyInputModeLookupTableShortcut,
			Type:      ibus.PROP_TYPE_NORMAL,
			Label:     dbus.MakeVariant(ibus.NewText("Phím tắt")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Keyboard Shortcuts")),
			Sensitive: true,
			Visible:   true,
			Icon:      "appointment",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
	)
}

func GetCharsetPropListByConfig(c *config.Config) *ibus.PropList {
	var charsetProperties []*ibus.Property
	charsetProperties = append(charsetProperties,
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyVnCharsetConvert,
			Type:      ibus.PROP_TYPE_NORMAL,
			Label:     dbus.MakeVariant(ibus.NewText("Chuyển mã online")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("")),
			Sensitive: true,
			Visible:   true,
			Symbol:    dbus.MakeVariant(ibus.NewText("C")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		IBusSeparator)
	for _, charset := range bamboo.GetCharsetNames() {
		var state = ibus.PROP_STATE_UNCHECKED
		if charset == c.OutputCharset {
			state = ibus.PROP_STATE_CHECKED
		}
		var imProp = &ibus.Property{
			Name:      "IBusProperty",
			Key:       "OutputCharset::" + charset,
			Type:      ibus.PROP_TYPE_RADIO,
			Label:     dbus.MakeVariant(ibus.NewText(charset)),
			Tooltip:   dbus.MakeVariant(ibus.NewText("OutputCharset: " + charset)),
			Sensitive: true,
			Visible:   true,
			State:     state,
			Symbol:    dbus.MakeVariant(ibus.NewText("U")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		}
		charsetProperties = append(charsetProperties, imProp)
	}
	return ibus.NewPropList(charsetProperties...)
}

func GetIMPropListByConfig(c *config.Config) *ibus.PropList {
	var imProperties []*ibus.Property
	imProperties = append(imProperties,
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyConfiguration,
			Type:      ibus.PROP_TYPE_NORMAL,
			Label:     dbus.MakeVariant(ibus.NewText("Tự định nghĩa kiểu gõ")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Tự định nghĩa kiểu gõ")),
			Sensitive: true,
			Visible:   true,
			Symbol:    dbus.MakeVariant(ibus.NewText("BC")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		IBusSeparator,
	)
	for im := range c.InputMethodDefinitions {
		var state = ibus.PROP_STATE_UNCHECKED
		if im == c.InputMethod {
			state = ibus.PROP_STATE_CHECKED
		}
		var imProp = &ibus.Property{
			Name:      "IBusProperty",
			Key:       im,
			Type:      ibus.PROP_TYPE_RADIO,
			Label:     dbus.MakeVariant(ibus.NewText(im)),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Kiểu gõ " + im)),
			Sensitive: true,
			Visible:   true,
			State:     state,
			Symbol:    dbus.MakeVariant(ibus.NewText("V")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		}
		imProperties = append(imProperties, imProp)
	}
	return ibus.NewPropList(imProperties...)
}

func GetMacroPropListByConfig(c *config.Config) *ibus.PropList {
	macroChecked := ibus.PROP_STATE_UNCHECKED
	autoCapitalizeMacro := ibus.PROP_STATE_UNCHECKED

	if c.IBflags&config.IBmacroEnabled != 0 {
		macroChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&config.IBautoCapitalizeMacro != 0 {
		autoCapitalizeMacro = ibus.PROP_STATE_CHECKED
	}
	return ibus.NewPropList(
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyMacroEnabled,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Bật gõ tắt")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Bật gõ tắt")),
			Sensitive: true,
			Visible:   true,
			State:     macroChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("M")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAutoCapitalizeMacro,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Tự động viết hoa")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Auto capitalize macro")),
			Sensitive: true,
			Visible:   true,
			State:     autoCapitalizeMacro,
			Symbol:    dbus.MakeVariant(ibus.NewText("C")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyMacroTable,
			Type:      ibus.PROP_TYPE_NORMAL,
			Label:     dbus.MakeVariant(ibus.NewText("Mở bảng gõ tắt")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Mở bảng gõ tắt")),
			Sensitive: true,
			Visible:   true,
			Symbol:    dbus.MakeVariant(ibus.NewText("O")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
	)
}

func GetSpellCheckingPropListByConfig(c *config.Config) *ibus.PropList {
	spellCheckByRules := ibus.PROP_STATE_UNCHECKED
	spellCheckByDicts := ibus.PROP_STATE_UNCHECKED

	// spelling
	spellingChecked := ibus.PROP_STATE_UNCHECKED
	if c.IBflags&config.IBspellCheckEnabled != 0 {
		spellingChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&config.IBspellCheckWithRules != 0 {
		spellCheckByRules = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&config.IBspellCheckWithDicts != 0 {
		spellCheckByDicts = ibus.PROP_STATE_CHECKED
	}
	return ibus.NewPropList(
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyEnableSpellCheck,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Bật kiểm tra chính tả")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("")),
			Sensitive: true,
			Visible:   true,
			State:     spellingChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("S")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		IBusSeparator,
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeySpellCheckByRules,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Sử dụng luật ghép vần")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Sử dụng luật ghép vần")),
			Sensitive: false,
			Visible:   true,
			State:     spellCheckByRules,
			Symbol:    dbus.MakeVariant(ibus.NewText("M")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeySpellCheckByDicts,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Sử dụng từ điển")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Sử dụng từ điển")),
			Sensitive: true,
			Visible:   true,
			State:     spellCheckByDicts,
			Symbol:    dbus.MakeVariant(ibus.NewText("O")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
	)
}

func GetOptionsPropListByConfig(c *config.Config) *ibus.PropList {
	// tone
	toneStdChecked := ibus.PROP_STATE_UNCHECKED
	toneFreeMarkingChecked := ibus.PROP_STATE_UNCHECKED
	preeditInvisibilityChecked := ibus.PROP_STATE_UNCHECKED
	x11FakeBackspaceChecked := ibus.PROP_STATE_UNCHECKED

	if c.Flags&bamboo.EstdToneStyle != 0 {
		toneStdChecked = ibus.PROP_STATE_CHECKED
	}
	if c.Flags&bamboo.EfreeToneMarking != 0 {
		toneFreeMarkingChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&config.IBnoUnderline != 0 {
		preeditInvisibilityChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&config.IBpreeditElimination != 0 {
		x11FakeBackspaceChecked = ibus.PROP_STATE_CHECKED
	}

	return ibus.NewPropList(
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyFreeToneMarking,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Bỏ dấu tự do")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Bỏ dấu tự do")),
			Sensitive: true,
			Visible:   true,
			State:     toneFreeMarkingChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("M")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyStdToneStyle,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Dấu thanh chuẩn")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Use òa, úy... (instead of oà, uý)")),
			Sensitive: true,
			Visible:   true,
			State:     toneStdChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("M")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyPreeditInvisibility,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Ẩn gạch chân")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Hide underline")),
			Sensitive: true,
			Visible:   true,
			State:     preeditInvisibilityChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("P")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyPreeditElimination,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Send key via ForwardKeyEvent")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Send key via ForwardKeyEvent")),
			Sensitive: false,
			Visible:   false,
			State:     x11FakeBackspaceChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("X")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
	)
}
