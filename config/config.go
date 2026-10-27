package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"

	"github.com/BambooEngine/bamboo-core"
	"ibus-bamboo/introspectors"
	"ibus-bamboo/types"
)

const (
	configDir        = "%s/.config/ibus-%s"
	configFile       = "%s/ibus-%s.config.json"
	mactabFile       = "%s/ibus-%s.macro.text"
	sampleMactabFile = "data/macro.tpl.txt"
)

func GetConfigDir(ngName string) string {
	u, err := user.Current()
	if err == nil {
		return fmt.Sprintf(configDir, u.HomeDir, "bamboo")
	}
	return fmt.Sprintf(configDir, "~", "bamboo")
}

func GetMacroPath(engineName string) string {
	return fmt.Sprintf(mactabFile, GetConfigDir(engineName), engineName)
}

func GetConfigPath(engineName string) string {
	return fmt.Sprintf(configFile, GetConfigDir(engineName), engineName)
}

func DefaultCfg() types.Config {
	return types.Config{
		InputMethod:            "Telex",
		OutputCharset:          "Unicode",
		InputMethodDefinitions: bamboo.GetInputMethodDefinitions(),
		Flags:                  bamboo.EstdFlags,
		IBflags:                IBstdFlags,
		Shortcuts:              [10]uint32{1, 126, 0, 0, 0, 0, 0, 0, 5, 117},
		DefaultInputMode:       PreeditIM,
		InputModeMapping:       map[string]int{},
		Introspector:           introspectors.IntrospectorAuto,
		IntrospectorCommand:    "",
	}
}

func LoadConfig(engineName string) *types.Config {
	var c = DefaultCfg()
	if engineName == "bamboous" {
		c.DefaultInputMode = UsIM
		c.IBflags = IBUsStdFlags
		return &c
	}

	var configPath = GetConfigPath(engineName)
	data, err := ioutil.ReadFile(configPath)
	if err == nil {
		json.Unmarshal(data, &c)
	}

	return &c
}

func SaveConfig(c *types.Config, engineName string) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf(configFile, GetConfigDir(engineName), engineName), data, 0644)
	if err != nil {
		log.Println(err)
	}

}
