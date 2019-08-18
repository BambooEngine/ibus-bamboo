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

const (
	IBUS_SHIFT_MASK   = 1 << 0
	IBUS_LOCK_MASK    = 1 << 1
	IBUS_CONTROL_MASK = 1 << 2
	IBUS_MOD1_MASK    = 1 << 3
	IBUS_MOD2_MASK    = 1 << 4
	IBUS_MOD3_MASK    = 1 << 5
	IBUS_MOD4_MASK    = 1 << 6
	IBUS_MOD5_MASK    = 1 << 7
	IBUS_BUTTON1_MASK = 1 << 8
	IBUS_BUTTON2_MASK = 1 << 9
	IBUS_BUTTON3_MASK = 1 << 10
	IBUS_BUTTON4_MASK = 1 << 11
	IBUS_BUTTON5_MASK = 1 << 12

	/* The next few modifiers are used by XKB so we skip to the end.
	 * Bits 15 - 23 are currently unused. Bit 29 is used internally.
	 */

	/* ibus mask */
	IBUS_HANDLED_MASK = 1 << 24
	IBUS_FORWARD_MASK = 1 << 25
	IBUS_IGNORED_MASK = IBUS_FORWARD_MASK

	IBUS_SUPER_MASK = 1 << 26
	IBUS_HYPER_MASK = 1 << 27
	IBUS_META_MASK  = 1 << 28

	IBUS_RELEASE_MASK = 1 << 30

	IBUS_MODIFIER_MASK = 0x5f001fff
)

const (
	//IBusCapability
	IBUS_CAP_PREEDIT_TEXT = 1 << 0 //UI is capable to show pre-edit text.
	//IBUS_CAP_AUXILIARY_TEXT   = 1 << 1 //UI is capable to show auxiliary text.
	//IBUS_CAP_LOOKUP_TABLE     = 1 << 2 //UI is capable to show the lookup table.
	//IBUS_CAP_FOCUS            = 1 << 3 //UI is capable to get focus.
	//IBUS_CAP_PROPERTY         = 1 << 4 //UI is capable to have property.
	IBUS_CAP_SURROUNDING_TEXT = 1 << 5 //Client can provide surround text, or IME can handle surround text.
)
const (
	IBUS_Colon            = 0x03a
	IBUS_Left             = 0xFF51
	IBUS_Up               = 0xFF52
	IBUS_Right            = 0xFF53
	IBUS_Down             = 0xFF54
	IBUS_Page_Up          = 0xFF55
	IBUS_Page_Down        = 0xFF56
	IBUS_BackSpace        = 0xff08
	IBUS_Return           = 0xff0d
	IBUS_Escape           = 0xff1b
	IBUS_Shift_L          = 0xffe1
	IBUS_Shift_R          = 0xffe2
	IBUS_Space            = 0x020
	IBUS_TILDE            = 0x007e
	IBUS_GRAVE            = 0x0060
	IBUS_Insert           = 0xff63
	IBUS_Deadkey_Currency = 0xfe6f
	IBUS_Caps_Lock        = 0xffe5
	IBUS_OpenLookupTable  = IBUS_TILDE
	IBUS_OpenEmojiTable   = IBUS_Colon
)

const (
	IBUS_ORIENTATION_HORIZONTAL = 0
	IBUS_ORIENTATION_VERTICAL   = 1
	IBUS_ORIENTATION_SYSTEM     = 2
)
