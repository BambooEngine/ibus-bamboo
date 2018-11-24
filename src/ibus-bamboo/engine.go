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
	"fmt"
	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
	"log"
	"os/exec"
	"sync"
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
	macroTable    *MacroTable
	vnSeq         string
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
		if e.config.IBflags&IBautoNonVnRestore == 0 {
			e.vnSeq = e.getVnSeq()
			e.preediter.Reset()
			var vnRunes = []rune(e.vnSeq)
			if len(vnRunes) == 0 {
				return false, nil
			}
			if len(vnRunes) >= 1 {
				vnRunes = vnRunes[:len(vnRunes)-1]
				e.vnSeq = string(vnRunes)
			} else {
				e.vnSeq = ""
			}
			e.updatePreedit()
			return true, nil
		}
		if rawKeyLen > 0 {
			e.preediter.RemoveLastChar()
			e.updatePreedit()
			return true, nil
		} else {
			return false, nil
		}
	}

	if keyVal == IBUS_space || keyVal == IBUS_KP_Space {
		e.ignorePreedit = false
		e.commitPreedit(0)
		return false, nil
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
	fmt.Printf("keyCode 0x%04x keyval 0x%04x | %c\n", keyCode, keyVal, rune(keyVal))

	if (keyVal >= 'a' && keyVal <= 'z') ||
		(keyVal >= 'A' && keyVal <= 'Z') ||
		(keyVal >= '0' && keyVal <= '9') ||
		(inKeyMap(e.preediter.GetInputMethod().Keys, rune(keyVal))) {
		var keyRune = rune(keyVal)
		if state&IBUS_LOCK_MASK != 0 {
			keyRune = toUpper(keyRune)
		}
		if e.ignorePreedit {
			return false, nil
		}
		if e.config.IBflags&IBautoNonVnRestore == 0 {
			var vnSeqTmp = e.preediter.GetProcessedString(bamboo.VietnameseMode)
			e.preediter.ProcessChar(keyRune, bamboo.VietnameseMode)
			if !e.preediter.IsSpellingLikelyCorrect(bamboo.NoTone) &&
				!inKeyMap(e.preediter.GetInputMethod().SuperKeys, keyRune) &&
				!inKeyMap(e.preediter.GetInputMethod().ToneKeys, keyRune) {
				e.vnSeq += vnSeqTmp
				e.preediter.Reset()
				e.preediter.ProcessChar(keyRune, bamboo.VietnameseMode)
			}
			log.Println(e.preediter.GetProcessedString(bamboo.VietnameseMode))
			e.updatePreedit()
			return true, nil
		}
		e.preediter.ProcessChar(keyRune, e.getMode())
		if e.config.IBflags&IBfastCommitEnabled != 0 && !e.preediter.IsSpellingLikelyCorrect(bamboo.NoTone) {
			e.ignorePreedit = true
			e.commitPreedit(0)
			e.preediter.Reset()
			return true, nil
		}
		e.updatePreedit()
		return true, nil
	} else {
		e.commitPreedit(keyVal)
		//forward lastKey
		e.ForwardKeyEvent(keyVal, keyCode, state)
		return true, nil
	}
	return false, nil
}

func (e *IBusBambooEngine) FocusIn() *dbus.Error {
	e.Lock()
	defer e.Unlock()

	e.RegisterProperties(e.propList)
	e.HidePreeditText()
	e.preediter.Reset()
	fmt.Print("FocusIn.")

	return nil
}

func (e *IBusBambooEngine) FocusOut() *dbus.Error {
	fmt.Print("FocusOut.")
	return nil
}

func (e *IBusBambooEngine) Reset() *dbus.Error {
	fmt.Print("Reset.")
	return nil
}

func (e *IBusBambooEngine) Enable() *dbus.Error {
	fmt.Print("Enable.")
	mouseCaptureInit()
	return nil
}

func (e *IBusBambooEngine) Disable() *dbus.Error {
	fmt.Print("Disable.")
	mouseCaptureExit()
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
	if propName == PropKeyAbout {
		exec.Command("xdg-open", HomePage).Start()
		return nil
	}
	if propName == PropKeyVnConvert {
		exec.Command("xdg-open", VnConvertPage).Start()
		return nil
	}
	if propName == PropKeyMacroTable {
		OpenMactabFile(EngineName)
	}

	turnSpellChecking := func(on bool) {
		if on {
			e.config.IBflags |= IBspellCheckEnabled
			e.config.IBflags |= IBautoNonVnRestore
		} else {
			e.config.IBflags &= ^IBspellCheckEnabled
			e.config.IBflags &= ^IBautoNonVnRestore
			e.config.IBflags &= ^IBfastCommitEnabled
		}
	}

	turnSpellCheckByRules := func(on bool) {
		turnSpellChecking(true)
		if on {
			e.config.IBflags |= IBspellCheckingByRules
			e.config.IBflags &= ^IBspellCheckingByDicts
			e.config.IBflags |= IBddFreeStyle
		} else {
			e.config.IBflags |= IBspellCheckingByDicts
			e.config.IBflags &= ^IBspellCheckingByRules
			e.config.IBflags &= ^IBddFreeStyle
			e.config.IBflags &= ^IBfastCommitEnabled
		}
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
			turnSpellChecking(true)
		} else {
			turnSpellChecking(false)
		}
	}
	if propName == PropKeySpellCheckingByRules {
		if propState == ibus.PROP_STATE_CHECKED {
			turnSpellCheckByRules(true)
		}
	}
	if propName == PropKeySpellCheckingByDicts {
		if propState == ibus.PROP_STATE_CHECKED {
			turnSpellCheckByRules(false)
		}
	}
	if propName == PropKeyFastCommit {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBfastCommitEnabled
			turnSpellCheckByRules(true)
		} else {
			e.config.IBflags &= ^IBfastCommitEnabled
		}
	}
	if propName == PropKeyMacroEnabled {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBmarcoEnabled
			e.config.IBflags &= ^IBfastCommitEnabled
			e.macroTable.Enable()
		} else {
			e.config.IBflags &= ^IBmarcoEnabled
			e.macroTable.Disable()
		}
	}
	if propName == PropKeyPreeditInvisibility {
		if propState == ibus.PROP_STATE_CHECKED {
			e.config.IBflags |= IBpreeditInvisibility
		} else {
			e.config.IBflags &= ^IBpreeditInvisibility
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
