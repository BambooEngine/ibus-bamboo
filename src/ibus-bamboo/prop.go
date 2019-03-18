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
	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

const (
	PropKeyAbout                       = "about"
	PropKeyStdToneStyle                = "tone_std_style"
	PropKeyFreeToneMarking             = "tone_free_marking"
	PropKeySpellingChecking            = "spelling_checking"
	PropKeySpellCheckingByRules        = "spelling_checking_by_rules"
	PropKeySpellCheckingByDicts        = "spelling_checking_by_dicts"
	PropKeyPreeditInvisibility         = "preedit_hiding"
	PropKeyVnConvert                   = "vn_convert"
	PropKeyAutoCommitWithVnNotMatch    = "AutoCommitWithSpellChecking"
	PropKeyAutoCommitWithVnFullMatch   = "AutoCommitWithVnFullMatch"
	PropKeyAutoCommitWithVnWordBreak   = "AutoCommitWithVnFC"
	PropKeyAutoCommitWithMouseMovement = "AutoCommitWithMouseMovement"
	PropKeyAutoCommitWithDelay         = "AutoCommitWithDelay"
	PropKeyMacroEnabled                = "macro_enabled"
	PropKeyMacroTable                  = "macro_table"
)

var runMode = ""

func GetPropListByConfig(c *Config) *ibus.PropList {
	return ibus.NewPropList(
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAbout,
			Type:      ibus.PROP_TYPE_NORMAL,
			Label:     dbus.MakeVariant(ibus.NewText("IBus " + EngineName + " " + Version + runMode)),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Mở trang chủ")),
			Sensitive: true,
			Visible:   true,
			Icon:      "gtk-home",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_SEPARATOR,
			Label:     dbus.MakeVariant(ibus.NewText("")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("")),
			Sensitive: true,
			Visible:   true,
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
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
			Label:     dbus.MakeVariant(ibus.NewText("Tự động gửi buffer")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Tự động gửi buffer")),
			Sensitive: true,
			Visible:   true,
			Icon:      "appointment",
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(GetAutoCommitPropListByConfig(c)),
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
	)
}

func GetCharsetPropListByConfig(c *Config) *ibus.PropList {
	var charsetProperties []*ibus.Property
	charsetProperties = append(charsetProperties,
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyVnConvert,
			Type:      ibus.PROP_TYPE_NORMAL,
			Label:     dbus.MakeVariant(ibus.NewText("Chuyển mã online")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("")),
			Sensitive: true,
			Visible:   true,
			Symbol:    dbus.MakeVariant(ibus.NewText("C")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_SEPARATOR,
			Label:     dbus.MakeVariant(ibus.NewText("")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("")),
			Sensitive: true,
			Visible:   true,
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		})
	for _, charset := range bamboo.GetCharsetNames() {
		var state = ibus.PROP_STATE_UNCHECKED
		if charset == c.Charset {
			state = ibus.PROP_STATE_CHECKED
		}
		var imProp = &ibus.Property{
			Name:      "IBusProperty",
			Key:       "Charset-" + charset,
			Type:      ibus.PROP_TYPE_RADIO,
			Label:     dbus.MakeVariant(ibus.NewText(charset)),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Charset: " + charset)),
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

func GetIMPropListByConfig(c *Config) *ibus.PropList {
	var imProperties []*ibus.Property
	for im := range bamboo.InputMethods {
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

func GetMacroPropListByConfig(c *Config) *ibus.PropList {
	macroChecked := ibus.PROP_STATE_UNCHECKED

	if c.IBflags&IBmarcoEnabled != 0 {
		macroChecked = ibus.PROP_STATE_CHECKED
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

func GetSpellCheckingPropListByConfig(c *Config) *ibus.PropList {
	spellCheckByRules := ibus.PROP_STATE_UNCHECKED
	spellCheckByDicts := ibus.PROP_STATE_UNCHECKED

	// spelling
	spellingChecked := ibus.PROP_STATE_UNCHECKED
	if c.IBflags&IBspellChecking != 0 {
		spellingChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&IBspellCheckingWithRules != 0 {
		spellCheckByRules = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&IBspellCheckingWithDicts != 0 {
		spellCheckByDicts = ibus.PROP_STATE_CHECKED
	}
	return ibus.NewPropList(
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeySpellingChecking,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Bật kiểm tra chính tả")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("")),
			Sensitive: true,
			Visible:   true,
			State:     spellingChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("S")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       "-",
			Type:      ibus.PROP_TYPE_SEPARATOR,
			Label:     dbus.MakeVariant(ibus.NewText("")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("")),
			Sensitive: true,
			Visible:   true,
			Symbol:    dbus.MakeVariant(ibus.NewText("")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeySpellCheckingByRules,
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
			Key:       PropKeySpellCheckingByDicts,
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

func GetOptionsPropListByConfig(c *Config) *ibus.PropList {
	// tone
	toneStdChecked := ibus.PROP_STATE_UNCHECKED
	toneFreeMarkingChecked := ibus.PROP_STATE_UNCHECKED
	preeditInvisibilityChecked := ibus.PROP_STATE_UNCHECKED

	if c.Flags&bamboo.EstdToneStyle != 0 {
		toneStdChecked = ibus.PROP_STATE_CHECKED
	}
	if c.Flags&bamboo.EfreeToneMarking != 0 {
		toneFreeMarkingChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&IBpreeditInvisibility != 0 {
		preeditInvisibilityChecked = ibus.PROP_STATE_CHECKED
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
			Tooltip:   dbus.MakeVariant(ibus.NewText("Ẩn gạch chân")),
			Sensitive: true,
			Visible:   true,
			State:     preeditInvisibilityChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("P")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
	)
}

func GetAutoCommitPropListByConfig(c *Config) *ibus.PropList {
	mouseMovementChecked := ibus.PROP_STATE_UNCHECKED
	vnFullMatchChecked := ibus.PROP_STATE_UNCHECKED
	vnWordBreakChecked := ibus.PROP_STATE_UNCHECKED
	vnNotMatchChecked := ibus.PROP_STATE_UNCHECKED
	delayChecked := ibus.PROP_STATE_UNCHECKED

	if c.IBflags&IBautoCommitWithMouseMovement != 0 {
		mouseMovementChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&IBautoCommitWithVnFullMatch != 0 {
		vnFullMatchChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&IBautoCommitWithVnNotMatch != 0 {
		vnNotMatchChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&IBautoCommitWithVnWordBreak != 0 {
		vnWordBreakChecked = ibus.PROP_STATE_CHECKED
	}
	if c.IBflags&IBautoCommitWithDelay != 0 {
		delayChecked = ibus.PROP_STATE_CHECKED
	}

	return ibus.NewPropList(
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAutoCommitWithVnNotMatch,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Sai chính tả")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Invalid words")),
			Sensitive: true,
			Visible:   true,
			State:     vnNotMatchChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("M")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAutoCommitWithVnFullMatch,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Hoàn thành một từ")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Finish a word")),
			Sensitive: false,
			Visible:   false,
			State:     vnFullMatchChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("M")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAutoCommitWithVnWordBreak,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Phụ âm đầu")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("first consonants")),
			Sensitive: true,
			Visible:   false,
			State:     vnWordBreakChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("P")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAutoCommitWithDelay,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Sau 3 giây")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("After 3s of inactive")),
			Sensitive: true,
			Visible:   true,
			State:     delayChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("F")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
		&ibus.Property{
			Name:      "IBusProperty",
			Key:       PropKeyAutoCommitWithMouseMovement,
			Type:      ibus.PROP_TYPE_TOGGLE,
			Label:     dbus.MakeVariant(ibus.NewText("Click hoặc di chuyển chuột")),
			Tooltip:   dbus.MakeVariant(ibus.NewText("Click or move mouse")),
			Sensitive: true,
			Visible:   true,
			State:     mouseMovementChecked,
			Symbol:    dbus.MakeVariant(ibus.NewText("F")),
			SubProps:  dbus.MakeVariant(*ibus.NewPropList()),
		},
	)
}
