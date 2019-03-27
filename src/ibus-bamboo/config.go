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
	"encoding/json"
	"fmt"
	"github.com/BambooEngine/bamboo-core"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

const (
	configFile       = "%s/.config/ibus/ibus-%s.config.json"
	mactabFile       = "%s/.config/ibus/ibus-%s.macro.text"
	sampleMactabFile = "data/macro.tpl.txt"
)

const (
	IBautoCommitWithVnNotMatch uint = 1 << iota
	IBmarcoEnabled
	IBautoCommitWithVnFullMatch
	IBautoCommitWithVnWordBreak
	IBspellChecking
	IBautoNonVnRestore
	IBddFreeStyle
	IBpreeditInvisibility
	IBspellCheckingWithRules
	IBspellCheckingWithDicts
	IBautoCommitWithDelay
	IBautoCommitWithMouseMovement
	IBemojiDisabled
	IBstdFlags = IBspellChecking | IBspellCheckingWithRules | IBautoNonVnRestore | IBddFreeStyle |
		IBpreeditInvisibility | IBautoCommitWithMouseMovement | IBemojiDisabled
)

type Config struct {
	InputMethod              string
	InputMethodDefinitions   map[string]bamboo.InputMethodDefinition
	Charset                  string
	Flags                    uint
	IBflags                  uint
	AutoCommitAfter          int64
	ExceptedWhiteList        []string
	PreeditWhiteList         []string
	X11ClipboardWhiteList    []string
	ForwardKeyWhiteList      []string
	SurroundingTextWhiteList []string
}

func getBambooConfigurationPath(engineName string) string {
	u, err := user.Current()
	if err == nil {
		return fmt.Sprintf(configFile, u.HomeDir, engineName)
	}
	return ""
}

func LoadConfig(engineName string) *Config {
	var c = Config{
		InputMethod:              "Telex",
		Charset:                  "Unicode",
		InputMethodDefinitions:   bamboo.InputMethodDefinitions,
		Flags:                    bamboo.EstdFlags,
		IBflags:                  IBstdFlags,
		AutoCommitAfter:          3000,
		ExceptedWhiteList:        nil,
		PreeditWhiteList:         nil,
		X11ClipboardWhiteList:    nil,
		ForwardKeyWhiteList:      nil,
		SurroundingTextWhiteList: nil,
	}

	data, err := ioutil.ReadFile(getBambooConfigurationPath(engineName))
	if err == nil {
		json.Unmarshal(data, &c)
	}

	return &c
}

func SaveConfig(c *Config, engineName string) {
	u, err := user.Current()
	if err != nil {
		return
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf(configFile, u.HomeDir, engineName), data, 0644)
	if err != nil {
		log.Println(err)
	}

}

func getEngineSubFile(fileName string) string {
	if _, err := os.Stat(fileName); err == nil {
		if absPath, err := filepath.Abs(fileName); err == nil {
			return absPath
		}
	}

	return filepath.Join(filepath.Dir(os.Args[0]), fileName)
}

func getMactabFile(engineName string) string {
	u, err := user.Current()
	if err != nil {
		return fmt.Sprintf(mactabFile, "~", engineName)
	}

	return fmt.Sprintf(mactabFile, u.HomeDir, engineName)
}

func OpenMactabFile(engineName string) {
	efPath := getMactabFile(engineName)
	if _, err := os.Stat(efPath); os.IsNotExist(err) {
		sampleFile := getEngineSubFile(sampleMactabFile)
		sample, err := ioutil.ReadFile(sampleFile)
		log.Println(err)
		ioutil.WriteFile(efPath, sample, 0644)
	}

	exec.Command("xdg-open", efPath).Start()
}
