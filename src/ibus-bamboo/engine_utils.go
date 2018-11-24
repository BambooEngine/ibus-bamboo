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
	"strings"
	"time"
)

func IBusBambooEngineCreator(conn *dbus.Conn, engineName string) dbus.ObjectPath {
	objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/bamboo/%d", time.Now().UnixNano()))

	var config = LoadConfig(engineName)

	engine := &IBusBambooEngine{
		Engine:     ibus.BaseEngine(conn, objectPath),
		preediter:  bamboo.NewEngine(config.InputMethod, config.Flags),
		engineName: engineName,
		config:     config,
		propList:   GetPropListByConfig(config),
		macroTable: NewMacroTable(),
	}
	ibus.PublishEngine(conn, objectPath, engine)

	if config.IBflags&IBmarcoEnabled!=0{
		engine.macroTable.Enable()
	}
	go engine.startAutoCommit()

	onMouseMove = func() {
		engine.Lock()
		defer engine.Unlock()
		var rawKeyLen = engine.getRawKeyLen()
		if rawKeyLen > 0 {
			engine.commitPreedit(0)
		}
	}

	return objectPath
}

var preeditUpdateChan = make(chan uint32)

func (e *IBusBambooEngine) startAutoCommit() {
	for {
		var timeout = e.config.AutoCommitAfter
		if e.config.IBflags&IBfastCommitEnabled == 0 {
			timeout = 10 * e.config.AutoCommitAfter
		}
		select {
		case <-preeditUpdateChan:
			break
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			var rawKeyLen = e.getRawKeyLen()
			if e.config.IBflags&IBautoCommitEnabled != 0 && rawKeyLen > 0 {
				e.commitPreedit(0)
			}
		}
	}
}

func (e *IBusBambooEngine) getRawKeyLen() int {
	return len(e.preediter.GetProcessedString(bamboo.EnglishMode))
}

func (e *IBusBambooEngine) updatePreedit() {
	var processedStr = e.getPreeditString()
	var preeditLen = uint32(len([]rune(processedStr)))
	var ibusText = ibus.NewText(processedStr)

	if e.config.IBflags&IBpreeditInvisibility != 0 {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_NONE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	} else {
		ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, preeditLen)
	}

	e.UpdatePreeditTextWithMode(ibusText, preeditLen, true, ibus.IBUS_ENGINE_PREEDIT_COMMIT)
	if preeditLen == 0 {
		e.HidePreeditText()
		e.preediter.Reset()
	}
	mouseCaptureUnlock()

	preeditUpdateChan <- 0
}

func (e *IBusBambooEngine) shouldFallbackToEnglish() bool {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return false
	}
	var vnSeq = e.preediter.GetProcessedString(bamboo.VietnameseMode)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return false
	}
	if e.config.IBflags&IBmarcoEnabled!=0 && e.macroTable.HasKey(vnSeq) {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && (vnRunes[len(vnRunes)-1] == 'd' || strings.ContainsRune(vnSeq, 'đ')) {
		return false
	}
	if e.preediter.IsSpellingCorrect(bamboo.NoTone) {
		return false
	}
	if e.preediter.IsSpellingLikelyCorrect(bamboo.NoTone) {
		return false
	}
	return true
}

func (e *IBusBambooEngine) mustFallbackToEnglish() bool {
	if e.config.IBflags&IBautoNonVnRestore == 0 {
		return false
	}
	var vnSeq = e.preediter.GetProcessedString(bamboo.VietnameseMode)
	var vnRunes = []rune(vnSeq)
	if len(vnRunes) == 0 {
		return false
	}
	// we want to allow dd even in non-vn sequence, because dd is used a lot in abbreviation
	if e.config.IBflags&IBddFreeStyle != 0 && strings.ContainsRune(vnSeq, 'đ') {
		return false
	}
	if e.preediter.IsSpellingCorrect(bamboo.NoTone) {
		return false
	}
	return true
}

func (e *IBusBambooEngine) getCommitString() string {
	var processedStr string
	processedStr = e.preediter.GetProcessedString(bamboo.VietnameseMode)
	if e.config.IBflags&IBmarcoEnabled!=0 && e.macroTable.HasKey(processedStr) {
		return e.macroTable.GetText(processedStr) + " "
	}
	if e.mustFallbackToEnglish() {
		processedStr = e.preediter.GetProcessedString(bamboo.EnglishMode)
		return processedStr
	}
	processedStr = bamboo.Encode(e.config.Charset, processedStr)
	return processedStr
}

func (e *IBusBambooEngine) getPreeditString() string {
	if e.shouldFallbackToEnglish() {
		return e.preediter.GetProcessedString(bamboo.EnglishMode)
	}
	return e.preediter.GetProcessedString(bamboo.VietnameseMode)
}

func (e *IBusBambooEngine) getMode() bamboo.Mode {
	if e.shouldFallbackToEnglish() {
		return bamboo.EnglishMode
	}
	return bamboo.VietnameseMode
}

func (e *IBusBambooEngine) commitPreedit(lastKey uint32) {
	var commitStr string
	commitStr += e.getCommitString()
	e.preediter.Reset()

	for _, chr := range []rune(commitStr) {
		e.CommitText(ibus.NewText(string(chr)))
	}
	//e.CommitText(ibus.NewText(commitStr))

	e.HidePreeditText()
}

func (e *IBusBambooEngine) hasMacroKey(key string) bool {
	return e.macroTable.GetText(key) != ""
}
func (e *IBusBambooEngine) commitMacroText() {
	var commitStr string
	var key = e.preediter.GetProcessedString(bamboo.VietnameseMode)
	commitStr = e.macroTable.GetText(key) + " "
	e.preediter.Reset()

	for _, chr := range []rune(commitStr) {
		e.CommitText(ibus.NewText(string(chr)))
	}
	//e.CommitText(ibus.NewText(commitStr))

	e.HidePreeditText()
}