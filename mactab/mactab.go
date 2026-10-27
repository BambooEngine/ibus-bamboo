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

package mactab

import (
	"bufio"
	"ibus-bamboo/config"
	"ibus-bamboo/types"
	"os"
	"strings"
	"time"
)

func NewMacroTable(autoCapitalizeMacro bool) *types.MacroTable {
	return &types.MacroTable{AutoCapitalizeMacro: autoCapitalizeMacro}
}

func LoadFromFile(e *types.MacroTable, macroFileName string) error {
	f, err := os.Open(macroFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	e.MTable = map[string]string{}
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
			if e.AutoCapitalizeMacro {
				key = strings.ToLower(key)
			}
			e.MTable[key] = strings.TrimSpace(list[1])
		}
	}
	return nil
}

func Reload(e *types.MacroTable, engineName string, autoCapitalizeMacro bool) {
	e.AutoCapitalizeMacro = autoCapitalizeMacro
	Enable(e, engineName)
}

func GetText(e *types.MacroTable, key string) string {
	if e.AutoCapitalizeMacro {
		key = strings.ToLower(key)
	}
	return e.MTable[key]
}

func HasKey(e *types.MacroTable, key string) bool {
	if e.AutoCapitalizeMacro {
		key = strings.ToLower(key)
	}
	return e.MTable[key] != ""
}

func HasPrefix(e *types.MacroTable, key string) bool {
	if e.MTable[key] != "" {
		return true
	}
	for k := range e.MTable {
		if strings.HasPrefix(k, key) {
			return true
		}
	}
	return false
}

func Enable(e *types.MacroTable, engineName string) {
	e.Enable = true

	go func() {
		modTime := time.Now()

		efPath := config.GetMacroPath(engineName)

		for e.Enable {
			if sta, _ := os.Stat(efPath); sta != nil {
				if newModeTime := sta.ModTime(); !newModeTime.Equal(modTime) {
					modTime = newModeTime
					LoadFromFile(e, efPath)
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

func Disable(e *types.MacroTable) {
	e.Enable = false
	e.MTable = map[string]string{}
}
