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
	"os/user"
	"path/filepath"
)

const (
	ComponentName = "org.freedesktop.IBus.bamboo"
	EngineName    = "Bamboo"
	HomePage      = "https://github.com/BambooEngine/ibus-bamboo"
	VnConvertPage = "https://tools.jcisio.com/vietuni/"

	DictVietnameseCm = "data/vietnamese.cm.dict"
	DictEmojiOne     = "data/emojione.json"
)

const (
	configDir        = "%s/.config/ibus-bamboo"
	configFile       = "%s/ibus-%s.config.json"
	mactabFile       = "%s/ibus-%s.macro.text"
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
	IBfakeBackspaceEnabled
	IBinputLookupTableDisabled
	IBstdFlags = IBspellChecking | IBspellCheckingWithRules | IBautoNonVnRestore | IBddFreeStyle |
		IBpreeditInvisibility | IBautoCommitWithMouseMovement | IBemojiDisabled | IBinputLookupTableDisabled
)

var DefaultBrowserList = []string{
	"Navigator:Firefox",
	"google-chrome:Google-chrome",
	"chromium-browser:Chromium-browser",
}

var DefaultPreeditList = []string{
	"google-chrome:Google-chrome",
}

type Config struct {
	InputMethod               string
	InputMethodDefinitions    map[string]bamboo.InputMethodDefinition
	OutputCharset             string
	Flags                     uint
	IBflags                   uint
	AutoCommitAfter           int64
	ExceptedList              []string
	PreeditWhiteList          []string
	X11ClipboardWhiteList     []string
	ForwardKeyWhiteList       []string
	DirectForwardKeyWhiteList []string
	SurroundingTextWhiteList  []string
}

func getConfigDir() string {
	u, err := user.Current()
	if err == nil {
		return fmt.Sprintf(configDir, u.HomeDir)
	}
	return fmt.Sprintf(configDir, "~")
}

func setupConfigDir() {
	if sta, err := os.Stat(getConfigDir()); err != nil || !sta.IsDir() {
		os.Mkdir(getConfigDir(), 0777)
	}
}

func getConfigPath(engineName string) string {
	return fmt.Sprintf(configFile, getConfigDir(), engineName)
}

func LoadConfig(engineName string) *Config {
	var c = Config{
		InputMethod:               "Telex",
		OutputCharset:             "Unicode",
		InputMethodDefinitions:    bamboo.InputMethodDefinitions,
		Flags:                     bamboo.EstdFlags,
		IBflags:                   IBstdFlags,
		AutoCommitAfter:           3000,
		ExceptedList:              nil,
		PreeditWhiteList:          DefaultPreeditList,
		X11ClipboardWhiteList:     nil,
		ForwardKeyWhiteList:       nil,
		DirectForwardKeyWhiteList: nil,
		SurroundingTextWhiteList:  nil,
	}

	data, err := ioutil.ReadFile(getConfigPath(engineName))
	if err == nil {
		json.Unmarshal(data, &c)
	}

	return &c
}

func SaveConfig(c *Config, engineName string) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf(configFile, getConfigDir(), engineName), data, 0644)
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
