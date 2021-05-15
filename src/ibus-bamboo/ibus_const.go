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
	IBusShiftMask   = 1 << 0
	IBusLockMask    = 1 << 1
	IBusControlMask = 1 << 2
	IBusMod1Mask    = 1 << 3

	/* The next few modifiers are used by XKB so we skip to the end.
	 * Bits 15 - 23 are currently unused. Bit 29 is used internally.
	 */

	/* ibus mask */
	IBusHandledMask = 1 << 24
	IBusForwardMask = 1 << 25
	IBusIgnoredMask = IBusForwardMask

	IBusSuperMask = 1 << 26
	IBusHyperMask = 1 << 27
	IBusMetaMask  = 1 << 28

	IBusReleaseMask = 1 << 30

	IBusModifierMask   = 0x5f001fff
	IBusDefaultModMask = IBusControlMask | IBusShiftMask | IBusMod1Mask | IBusSuperMask | IBusHyperMask | IBusMetaMask
)

const (
	//IBusCapability
	IBusCapPreeditText = 1 << 0 //UI is capable to show pre-edit text.
	//IBUS_CAP_AUXILIARY_TEXT   = 1 << 1 //UI is capable to show auxiliary text.
	//IBUS_CAP_LOOKUP_TABLE     = 1 << 2 //UI is capable to show the lookup table.
	//IBUS_CAP_FOCUS            = 1 << 3 //UI is capable to get focus.
	//IBUS_CAP_PROPERTY         = 1 << 4 //UI is capable to have property.
	IBusCapSurroundingText = 1 << 5 //Client can provide surround text, or IME can handle surround text.
)
const (
	XkBackspace = 0x16
	XkLeft      = 0x71
)
const (
	IBusTab             = 0xff09
	IBusEnd             = 0xff57
	IBusColon           = 0x03a
	IBusLeft            = 0xFF51
	IBusUp              = 0xFF52
	IBusRight           = 0xFF53
	IBusDown            = 0xFF54
	IBusPageUp          = 0xFF55
	IBusPageDown        = 0xFF56
	IBusBackSpace       = 0xff08
	IBusReturn          = 0xff0d
	IBusEscape          = 0xff1b
	IBusShiftL          = 0xffe1
	IBusShiftR          = 0xffe2
	IBusSpace           = 0x020
	IBusTilde           = 0x007e
	IBusGrave           = 0x0060
	IBusInsert          = 0xff63
	IBusCapsLock        = 0xffe5
	IBusOpenLookupTable = IBusTilde
	IBusOpenEmojiTable  = IBusColon
)

const (
	IBusOrientationHorizontal = 0
	IBusOrientationVertical   = 1
	IBusOrientationSystem     = 2
)
