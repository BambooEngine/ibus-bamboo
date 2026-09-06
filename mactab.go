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
	"bufio"
	"ibus-bamboo/config"
	"os"
	"strings"
	"sync"
	"time"
)

type MacroTable struct {
	sync.RWMutex
	enable              bool
	watching            bool
	autoCapitalizeMacro bool
	mTable              map[string]string
}

func NewMacroTable(autoCapitalizeMacro bool) *MacroTable {
	return &MacroTable{autoCapitalizeMacro: autoCapitalizeMacro}
}

func (e *MacroTable) LoadFromFile(macroFileName string) error {
	f, err := os.Open(macroFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	e.Lock()
	defer e.Unlock()
	e.mTable = map[string]string{}
	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			break
		}
		var s = strings.TrimSpace(string(line))
		if len(line) == 0 || strings.HasPrefix(s, ";") || strings.HasPrefix(s, "#") {
			continue
		}
		var list = strings.Split(s, ":")
		if len(list) == 2 {
			key := strings.TrimSpace(list[0])
			if e.autoCapitalizeMacro {
				key = strings.ToLower(key)
			}
			e.mTable[key] = strings.TrimSpace(list[1])
		}
	}
	return nil
}

func (e *MacroTable) Reload(engineName string, autoCapitalizeMacro bool) {
	e.Lock()
	e.autoCapitalizeMacro = autoCapitalizeMacro
	e.Unlock()
	e.Enable(engineName)
}

func (e *MacroTable) GetText(key string) string {
	e.RLock()
	defer e.RUnlock()
	if e.autoCapitalizeMacro {
		key = strings.ToLower(key)
	}
	return e.mTable[key]
}

func (e *MacroTable) HasKey(key string) bool {
	e.RLock()
	defer e.RUnlock()
	if e.autoCapitalizeMacro {
		key = strings.ToLower(key)
	}
	return e.mTable[key] != ""
}

func (e *MacroTable) HasPrefix(key string) bool {
	e.RLock()
	defer e.RUnlock()
	if e.mTable[key] != "" {
		return true
	}
	for k := range e.mTable {
		if strings.HasPrefix(k, key) {
			return true
		}
	}
	return false
}

func (e *MacroTable) isEnabled() bool {
	e.RLock()
	defer e.RUnlock()
	return e.enable
}

func (e *MacroTable) Enable(engineName string) {
	e.Lock()
	e.enable = true
	alreadyWatching := e.watching
	if !alreadyWatching {
		e.watching = true
	}
	e.Unlock()
	if alreadyWatching {
		return
	}

	go func() {
		defer func() {
			e.Lock()
			e.watching = false
			e.Unlock()
		}()
		modTime := time.Now()

		efPath := config.GetMacroPath(engineName)

		for e.isEnabled() {
			if sta, _ := os.Stat(efPath); sta != nil {
				if newModeTime := sta.ModTime(); !newModeTime.Equal(modTime) {
					modTime = newModeTime
					e.LoadFromFile(efPath)
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

func (e *MacroTable) Disable() {
	e.Lock()
	defer e.Unlock()
	e.enable = false
	e.mTable = map[string]string{}
}
