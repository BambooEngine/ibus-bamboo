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
	"encoding/json"
	"fmt"
	"github.com/BambooEngine/bamboo-core"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const (
	VnCaseAllSmall uint8 = 1 << iota
	VnCaseAllCapital
	VnCaseNoChange
)
const (
	HomePage           = "https://github.com/BambooEngine/ibus-bamboo"
	CharsetConvertPage = "https://tools.jcisio.com/vietuni/"

	DataDir          = "/usr/share/ibus-bamboo"
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
	preeditIM = iota + 1
	surroundingTextIM
	backspaceForwardingIM
	shiftLeftForwardingIM
	xTestFakeKeyEventIM
	forwardAsCommitIM
	usIM
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
	IBinputModeLookupTableEnabled
	IBautoCapitalizeMacro
	IBimQuickSwitchEnabled
	IBrestoreKeyStrokesEnabled
	IBmouseCapturing
	IBstdFlags = IBspellChecking | IBspellCheckingWithRules | IBautoNonVnRestore | IBddFreeStyle |
		IBemojiDisabled | IBinputModeLookupTableEnabled
)

var DefaultBrowserList = []string{
	"Navigator:Firefox",
	"google-chrome:Google-chrome",
	"chromium-browser:Chromium-browser",
}

type Config struct {
	InputMethod               string
	InputMethodDefinitions    map[string]bamboo.InputMethodDefinition
	OutputCharset             string
	Flags                     uint
	IBflags                   uint
	DefaultInputMode          int
	ExceptedList              []string
	PreeditWhiteList          []string
	X11ClipboardWhiteList     []string
	ForwardKeyWhiteList       []string
	SLForwardKeyWhiteList     []string
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
		InputMethodDefinitions:    bamboo.GetInputMethodDefinitions(),
		Flags:                     bamboo.EstdFlags,
		IBflags:                   IBstdFlags,
		DefaultInputMode:          preeditIM,
		ExceptedList:              nil,
		PreeditWhiteList:          nil,
		X11ClipboardWhiteList:     nil,
		ForwardKeyWhiteList:       nil,
		SLForwardKeyWhiteList:     nil,
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

func determineMacroCase(str string) uint8 {
	var chars = []rune(str)
	if unicode.IsLower(chars[0]) {
		return VnCaseAllSmall
	} else {
		for _, c := range chars[1:] {
			if unicode.IsLower(c) {
				return VnCaseNoChange
			}
		}
	}
	return VnCaseAllCapital
}

func toUpper(keyRune rune) rune {
	var upperSpecialKeys = map[rune]rune{
		'[': '{',
		']': '}',
	}

	if upperSpecialKey, found := upperSpecialKeys[keyRune]; found {
		keyRune = upperSpecialKey
	} else {
		keyRune = unicode.ToUpper(keyRune)
	}
	return keyRune
}

func inStringList(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func removeFromWhiteList(list []string, classes string) []string {
	var newList []string
	for _, cl := range list {
		if cl != classes {
			newList = append(newList, cl)
		}
	}
	return newList
}

func addToWhiteList(list []string, classes string) []string {
	for _, cl := range list {
		if cl == classes {
			return list
		}
	}
	return append(list, classes)
}

func getCharsetFromPropKey(str string) (string, bool) {
	var arr = strings.Split(str, "::")
	if len(arr) == 2 {
		return arr[1], true
	}
	return str, false
}

func isValidCharset(str string) bool {
	var charsets = bamboo.GetCharsetNames()
	for _, cs := range charsets {
		if cs == str {
			return true
		}
	}
	return false
}

type byString []string

func (s byString) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s byString) Len() int {
	return len(s)
}
func (s byString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sortStrings(list []string) []string {
	var strList = byString(list)
	sort.Sort(strList)
	return strList
}

func loadDictionary(dataFiles ...string) (map[string]bool, error) {
	var data = map[string]bool{}
	for _, dataFile := range dataFiles {
		f, err := os.Open(dataFile)
		if err != nil {
			return nil, err
		}
		rd := bufio.NewReader(f)
		for {
			line, _, err := rd.ReadLine()
			if err != nil {
				break
			}
			if len(line) == 0 {
				continue
			}
			var tmp = []byte(strings.ToLower(string(line)))
			data[string(tmp)] = true
			//bamboo.AddTrie(rootWordTrie, []rune(string(line)), false)
		}
		f.Close()
	}
	return data, nil
}

func isMovementKey(keyVal uint32) bool {
	var list = []uint32{IBUS_Left, IBUS_Right, IBUS_Up, IBUS_Down, IBUS_Page_Down, IBUS_Page_Up, IBUS_End, IBUS_KP_Down,
		IBUS_KP_End, IBUS_KP_Left, IBUS_KP_Next, IBUS_KP_Page_Down, IBUS_KP_Page_Up, IBUS_KP_Right, IBUS_KP_Up}
	for _, item := range list {
		if item == keyVal {
			return true
		}
	}
	return false
}

var vnSymMapping = map[rune]uint32{
	'Ạ': 0x1001ea0,
	'ạ': 0x1001ea1,
	'Ả': 0x1001ea2,
	'ả': 0x1001ea3,
	'Ấ': 0x1001ea4,
	'ấ': 0x1001ea5,
	'Ầ': 0x1001ea6,
	'ầ': 0x1001ea7,
	'Ẩ': 0x1001ea8,
	'ẩ': 0x1001ea9,
	'Ẫ': 0x1001eaa,
	'ẫ': 0x1001eab,
	'Ậ': 0x1001eac,
	'ậ': 0x1001ead,
	'Ắ': 0x1001eae,
	'ắ': 0x1001eaf,
	'Ằ': 0x1001eb0,
	'ằ': 0x1001eb1,
	'Ẳ': 0x1001eb2,
	'ẳ': 0x1001eb3,
	'Ẵ': 0x1001eb4,
	'ẵ': 0x1001eb5,
	'Ặ': 0x1001eb6,
	'ặ': 0x1001eb7,
	'Ẹ': 0x1001eb8,
	'ẹ': 0x1001eb9,
	'Ẻ': 0x1001eba,
	'ẻ': 0x1001ebb,
	'Ẽ': 0x1001ebc,
	'ẽ': 0x1001ebd,
	'Ế': 0x1001ebe,
	'ế': 0x1001ebf,
	'Ề': 0x1001ec0,
	'ề': 0x1001ec1,
	'Ể': 0x1001ec2,
	'ể': 0x1001ec3,
	'Ễ': 0x1001ec4,
	'ễ': 0x1001ec5,
	'Ệ': 0x1001ec6,
	'ệ': 0x1001ec7,
	'Ỉ': 0x1001ec8,
	'ỉ': 0x1001ec9,
	'Ị': 0x1001eca,
	'ị': 0x1001ecb,
	'Ọ': 0x1001ecc,
	'ọ': 0x1001ecd,
	'Ỏ': 0x1001ece,
	'ỏ': 0x1001ecf,
	'Ố': 0x1001ed0,
	'ố': 0x1001ed1,
	'Ồ': 0x1001ed2,
	'ồ': 0x1001ed3,
	'Ổ': 0x1001ed4,
	'ổ': 0x1001ed5,
	'Ỗ': 0x1001ed6,
	'ỗ': 0x1001ed7,
	'Ộ': 0x1001ed8,
	'ộ': 0x1001ed9,
	'Ớ': 0x1001eda,
	'ớ': 0x1001edb,
	'Ờ': 0x1001edc,
	'ờ': 0x1001edd,
	'Ở': 0x1001ede,
	'ở': 0x1001edf,
	'Ỡ': 0x1001ee0,
	'ỡ': 0x1001ee1,
	'Ợ': 0x1001ee2,
	'ợ': 0x1001ee3,
	'Ụ': 0x1001ee4,
	'ụ': 0x1001ee5,
	'Ủ': 0x1001ee6,
	'ủ': 0x1001ee7,
	'Ứ': 0x1001ee8,
	'ứ': 0x1001ee9,
	'Ừ': 0x1001eea,
	'ừ': 0x1001eeb,
	'Ử': 0x1001eec,
	'ử': 0x1001eed,
	'Ữ': 0x1001eee,
	'ữ': 0x1001eef,
	'Ự': 0x1001ef0,
	'ự': 0x1001ef1,
	'Ỵ': 0x1001ef4,
	'ỵ': 0x1001ef5,
	'Ỷ': 0x1001ef6,
	'ỷ': 0x1001ef7,
	'Ỹ': 0x1001ef8,
	'ỹ': 0x1001ef9,
	'Ơ': 0x10001a0,
	'ơ': 0x10001a1,
	'Ư': 0x10001af,
	'ư': 0x10001b0,
	'ă': 0x01e3,
	'Ă': 0x01c3,
	'Ỳ': 0x1001ef2,
	'ỳ': 0x1001ef3,
	'Đ': 0x01d0,
	'đ': 0x01f0,
	'Ĩ': 0x03a5,
	'ĩ': 0x03b5,
	'Ũ': 0x03dd,
	'ũ': 0x03fd,
}
