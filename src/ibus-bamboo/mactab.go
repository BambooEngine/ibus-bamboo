package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MACRO_VERSION_UTF8 uint = 1

type MacroTable struct {
	sync.RWMutex
	enable bool
	mTable map[string]string
}

func NewMacroTable() *MacroTable {
	return &MacroTable{}
}

//----------------------------------------------------------------------------
// Read header, if it's present in the file. Get the version of the file
// If header is absent, go back to the beginning of file and set version to 0
// Return false if reading failed.
//
// Header format: # DO NOT DELETE THIS LINE*** version=n ***
//----------------------------------------------------------------------------
func (e *MacroTable) ReadMacVersion(rd *bufio.Reader) (uint, error) {
	e.Lock()
	defer e.Unlock()
	line, _, err := rd.ReadLine()
	if err != nil || len(line) <= 0 {
		return 0, nil
	}
	var s = strings.TrimSpace(string(line))
	var reg = regexp.MustCompile(`DO NOT DELETE THIS LINE\*\*\*\sversion=(\d+)\s\*\*\*`)
	if reg.MatchString(s) {
		var matches = reg.FindStringSubmatch(s)
		var version, err = strconv.Atoi(matches[1])
		return uint(version), err
	}
	return 0, errors.New("invalid header")
}

//---------------------------------------------------------------
func (e *MacroTable) WriteHeader(f *os.File) {
	f.WriteString(fmt.Sprintf("# DO NOT DELETE THIS LINE*** version=%d ***\n", MACRO_VERSION_UTF8))
}

//---------------------------------------------------------------
func (e *MacroTable) LoadFromFile(macroFileName string) error {
	f, err := os.Open(macroFileName)
	defer f.Close()
	if err != nil {
		return err
	}
	e.mTable = map[string]string{}
	rd := bufio.NewReader(f)
	var version uint
	if version, err = e.ReadMacVersion(rd); err != nil {
		version = 0
	}
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
			e.mTable[list[0]] = list[1]
		}
	}
	// Convert old version
	if version != MACRO_VERSION_UTF8 {
		//e.WriteToFile(macroFileName)
	}
	return nil
}

//---------------------------------------------------------------
func (e *MacroTable) WriteToFile(macroFileName string) error {
	f, err := os.Open(macroFileName)
	defer f.Close()
	if err != nil {
		return err
	}
	e.WriteHeader(f)
	for key, text := range e.mTable {
		f.WriteString(fmt.Sprintf("%s:%s\n", key, text))
	}
	return nil
}

//---------------------------------------------------------------
func (e *MacroTable) GetText(key string) string {
	return e.mTable[key]
}

//---------------------------------------------------------------
func (e *MacroTable) HasKey(key string) bool {
	return e.mTable[key] != ""
}

//---------------------------------------------------------------
func (e *MacroTable) IncludeKey(key string) bool {
	if e.mTable[key] != "" {
		return true
	}
	for k, _ := range e.mTable {
		if strings.Contains(k, key) {
			return true
		}
	}
	return false
}

//---------------------------------------------------------------
func (e *MacroTable) Enable() {
	e.Lock()
	defer e.Unlock()
	e.enable = true

	go func() {
		cont := true
		modTime := time.Now()

		efPath := getMactabFile(EngineName)

		for cont {
			if sta, _ := os.Stat(efPath); sta != nil {
				if newModeTime := sta.ModTime(); !newModeTime.Equal(modTime) {
					modTime = newModeTime
					e.LoadFromFile(efPath)
				}
			}
			time.Sleep(time.Second)
			e.RLock()
			cont = e.enable
			e.RUnlock()
		}
	}()
}

//---------------------------------------------------------------
func (e *MacroTable) Disable() {
	e.Lock()
	defer e.Unlock()
	e.enable = false
	e.mTable = map[string]string{}
}
