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

//ibustypes — Generic types for IBus.
// http://ibus.github.io/docs/ibus-1.5/ibus-ibustypes.html

const (
	//IBusModifierType
	IBUS_SHIFT_MASK   = 1 << 0 //Shift is activated.
	IBUS_LOCK_MASK    = 1 << 1 //Cap Lock is locked.
	IBUS_CONTROL_MASK = 1 << 2 //Control key is activated.
	IBUS_MOD1_MASK    = 1 << 3 //Modifier 1 (Usually Alt_L (0x40), Alt_R (0x6c), Meta_L (0xcd)) activated.
	//IBUS_MOD2_MASK     = 1 << 4            //Modifier 2 (Usually Num_Lock (0x4d)) activated.
	//IBUS_MOD3_MASK     = 1 << 5            //Modifier 3 activated.
	//IBUS_MOD4_MASK     = 1 << 6            //Modifier 4 (Usually Super_L (0xce), Hyper_L (0xcf)) activated.
	//IBUS_MOD5_MASK     = 1 << 7            //Modifier 5 (ISO_Level3_Shift (0x5c), Mode_switch (0xcb)) activated.
	//IBUS_BUTTON1_MASK  = 1 << 8            //Mouse button 1 (left) is activated.
	//IBUS_BUTTON2_MASK  = 1 << 9            //Mouse button 2 (middle) is activated.
	//IBUS_BUTTON3_MASK  = 1 << 10           //Mouse button 3 (right) is activated.
	//IBUS_BUTTON4_MASK  = 1 << 11           //Mouse button 4 (scroll up) is activated.
	//IBUS_BUTTON5_MASK  = 1 << 12           //Mouse button 5 (scroll down) is activated.
	//IBUS_HANDLED_MASK  = 1 << 24           //Handled mask indicates the event has been handled by ibus.
	IBUS_FORWARD_MASK = 1 << 25           //Forward mask indicates the event has been forward from ibus.
	IBUS_IGNORED_MASK = IBUS_FORWARD_MASK //It is an alias of IBUS_FORWARD_MASK.
	IBUS_SUPER_MASK   = 1 << 26           //Super (Usually Win) key is activated.
	IBUS_HYPER_MASK   = 1 << 27           //Hyper key is activated.
	IBUS_META_MASK    = 1 << 28           //Meta key is activated.
	IBUS_RELEASE_MASK = 1 << 30           //Key is released.
	//IBUS_MODIFIER_MASK = 0x5f001fff        //Modifier mask for the all the masks above.
)

const (
	IBUS_INPUT_PURPOSE_FREE_FORM = iota
	IBUS_INPUT_PURPOSE_ALPHA
	IBUS_INPUT_PURPOSE_DIGITS
	IBUS_INPUT_PURPOSE_NUMBER
	IBUS_INPUT_PURPOSE_PHONE
	IBUS_INPUT_PURPOSE_URL
	IBUS_INPUT_PURPOSE_EMAIL
	IBUS_INPUT_PURPOSE_NAME
	IBUS_INPUT_PURPOSE_PASSWORD
	IBUS_INPUT_PURPOSE_PIN
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

// ibuskeysyms-compat
// http://ibus.github.io/docs/ibus-1.5/ibus-ibuskeysyms-compat.htm

const (
	IBUS_BackSpace        = 0xff08
	IBUS_Return           = 0xff0d
	IBUS_Escape           = 0xff1b
	IBUS_KP_Space         = 0xff80
	IBUS_KP_Enter         = 0xff8d
	IBUS_KP_Multiply      = 0xffaa
	IBUS_KP_Divide        = 0xffaf
	IBUS_KP_0             = 0xffb0
	IBUS_KP_9             = 0xffb9
	IBUS_Shift_L          = 0xffe1
	IBUS_Shift_R          = 0xffe2
	IBUS_space            = 0x020
	IBUS_0                = 0x030
	IBUS_TILDE            = 0x007e
	IBUS_GRAVE            = 0x0060
	IBUS_KEY_Insert       = 0xff63
	IBUS_Deadkey_Currency = 0xfe6f
	IBUS_Deadkey          = 96
	IBUS_OpenLookupTable  = IBUS_TILDE
)

const (
	IBUS_ORIENTATION_HORIZONTAL = 0
	IBUS_ORIENTATION_VERTICAL   = 1
	IBUS_ORIENTATION_SYSTEM     = 2
)
