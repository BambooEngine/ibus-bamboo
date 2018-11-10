/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Nguyen Cong Hoang <hoangnc.jp@gmail.com>
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
	"log"
	"os/exec"
	"runtime/debug"
	"sync"
	"unicode"
)

type IBusBambooEngine struct {
	sync.Mutex
	ibus.Engine
	preediter     bamboo.IEngine
	zeroLocation  bool
	engineName    string
	config        *Config
	propList      *ibus.PropList
	mode          bamboo.Mode
	ignorePreedit bool
}

/**
Implement IBus.Engine's process_key_event default signal handler.

Args:
	keyval - The keycode, transformed through a keymap, stays the
		same for every keyboard
	keycode - Keyboard-dependant key code
	modifiers - The state of IBus.ModifierType keys like
		Shift, Control, etc.
Return:
	True - if successfully process the keyevent, it won't be sent to X-server
	False - otherwise.

This function gets called whenever a key is pressed.
*/
func (e *IBusBambooEngine) ProcessKeyEvent(keyVal uint32, keyCode uint32, state uint32) (bool, *dbus.Error) {
	e.Lock()
	defer e.Unlock()
	var rawKeyLen = e.getRawKeyLen()

	if e.zeroLocation ||
		state&IBUS_RELEASE_MASK != 0 || //Ignore key-up event
		(state&IBUS_SHIFT_MASK == 0 && (keyVal == IBUS_Shift_L || keyVal == IBUS_Shift_R)) { //Ignore 1 shift key
		return false, nil
	}

	keyPressChan <- keyVal

	if state&IBUS_CONTROL_MASK != 0 ||
		state&IBUS_MOD1_MASK != 0 ||
		state&IBUS_IGNORED_MASK != 0 ||
		state&IBUS_SUPER_MASK != 0 ||
		state&IBUS_HYPER_MASK != 0 ||
		state&IBUS_META_MASK != 0 {
		e.ignorePreedit = false
		if rawKeyLen == 0 {
			return false, nil
		} else {
			//while typing, do not process control keys
			return true, nil
		}
	}

	if keyVal == IBUS_BackSpace {
		e.ignorePreedit = false
		if rawKeyLen > 0 {
			e.preediter.RemoveLastChar()
			e.updatePreedit()
			return true, nil
		}
	}

	if keyVal == IBUS_space || keyVal == IBUS_KP_Space {
		e.ignorePreedit = false
	}

	if keyVal == IBUS_Return || keyVal == IBUS_KP_Enter {
		e.ignorePreedit = false
		if rawKeyLen > 0 {
			e.commitPreedit(keyVal)
			e.ForwardKeyEvent(keyVal, keyCode, state)
			return true, nil
		} else {
			return false, nil
		}
	}

	if keyVal == IBUS_Escape {
		e.ignorePreedit = false
		if rawKeyLen > 0 {
			e.commitPreedit(keyVal)
			return true, nil
		}
		return false, nil
	}
	log.Printf("keyCode 0x%04x keyval 0x%04x | %c", keyCode, keyVal, rune(keyVal))

	if (keyVal >= 'a' && keyVal <= 'z') ||
		(keyVal >= 'A' && keyVal <= 'Z') ||
		(inKeyMap(e.preediter.GetInputMethod().Keys, rune(keyVal))) {
		var keyRune = rune(keyVal)
		if state&IBUS_LOCK_MASK != 0 {
			if upperSpecialKey, found := upperSpecialKeys[keyRune]; found {
				keyRune = upperSpecialKey
			} else {
				keyRune = unicode.ToUpper(keyRune)
			}
		}
		if e.ignorePreedit {
			return false, nil
		}
		e.preediter.ProcessChar(keyRune)
		if e.config.Flags&bamboo.EfastCommitting != 0 && !e.preediter.IsLikelySpellingCorrect(bamboo.NoTone) {
			e.ignorePreedit = true
			e.commitPreedit(0)
			e.preediter.Reset()
			return true, nil
		}
		e.updatePreedit()
		return true, nil
	} else {
		if rawKeyLen > 0 {
			if e.commitPreedit(keyVal) {
				//lastKey already appended to commit string
				return true, nil
			} else {
				//forward lastKey
				e.ForwardKeyEvent(keyVal, keyCode, state)
				return true, nil
			}
		}
		//pre-edit empty, just forward key
		return false, nil
	}
}

func (e *IBusBambooEngine) FocusIn() *dbus.Error {
	e.Lock()
	defer e.Unlock()

	e.RegisterProperties(e.propList)
	e.HidePreeditText()
	e.preediter.Reset()

	return nil
}

func (e *IBusBambooEngine) FocusOut() *dbus.Error {
	return nil
}

func (e *IBusBambooEngine) Reset() *dbus.Error {
	return nil
}

func (e *IBusBambooEngine) Enable() *dbus.Error {
	return nil
}

func (e *IBusBambooEngine) Disable() *dbus.Error {
	return nil
}

func (e *IBusBambooEngine) SetCapabilities(cap uint32) *dbus.Error {
	return nil
}

func (e *IBusBambooEngine) SetCursorLocation(x int32, y int32, w int32, h int32) *dbus.Error {
	e.zeroLocation = x == 0 && y == 0 && w == 0 && h == 0
	return nil
}

func (e *IBusBambooEngine) SetContentType(purpose uint32, hints uint32) *dbus.Error {
	return nil
}

//@method(in_signature="su")
func (e *IBusBambooEngine) PropertyActivate(propName string, propState uint32) *dbus.Error {
	debug.FreeOSMemory()

	if propName == PropKeyAbout {
		exec.Command("xdg-open", HomePage).Start()
		return nil
	}
	if propName == PropKeyVnConvert {
		exec.Command("xdg-open", VnConvertPage).Start()
		return nil
	}

	if propName == PropKeyStdToneStyle {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.Flags |= bamboo.EstdToneStyle
		} else {
			e.config.Flags &= ^bamboo.EstdToneStyle
		}
	}
	if propName == PropKeyFreeToneMarking {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.Flags |= bamboo.EfreeToneMarking
		} else {
			e.config.Flags &= ^bamboo.EfreeToneMarking
		}
	}
	if propName == PropKeySpellingChecking {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.Flags |= bamboo.EspellCheckEnabled
		} else {
			e.config.Flags &= ^bamboo.EspellCheckEnabled
		}
	}
	if propName == PropKeyFastCommitting {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.Flags |= bamboo.EfastCommitting
		} else {
			e.config.Flags &= ^bamboo.EfastCommitting
		}
	}
	if isValidCharset(getCharsetFromPropKey(propName)) && propState == ibus.PROP_STATE_CHECKED {
		e.config.Charset = getCharsetFromPropKey(propName)
	}
	if _, found := bamboo.InputMethods[propName]; found && propState == ibus.PROP_STATE_CHECKED {
		e.config.InputMethod = propName
	}
	SaveConfig(e.config, e.engineName)
	e.propList = GetPropListByConfig(e.config)
	e.preediter = bamboo.NewEngine(e.config.InputMethod, e.config.Flags)
	e.RegisterProperties(e.propList)
	return nil
}
