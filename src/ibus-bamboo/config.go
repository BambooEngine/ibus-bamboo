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
	"encoding/json"
	"fmt"
	"github.com/BambooEngine/bamboo-core"
	"io/ioutil"
	"log"
	"os/user"
)

const (
	configFile = "%s/.config/ibus/ibus-%s.config.json"
)

type Config struct {
	InputMethod string
	Charset     string
	Flags       uint
}

func LoadConfig(engineName string) *Config {
	var c = Config{
		InputMethod: "Telex 2",
		Charset:     "Unicode",
		Flags:       bamboo.EstdFlags,
	}

	u, err := user.Current()
	if err == nil {
		data, err := ioutil.ReadFile(fmt.Sprintf(configFile, u.HomeDir, engineName))
		if err == nil {
			json.Unmarshal(data, &c)
		}
	}

	return &c
}

func SaveConfig(c *Config, engineName string) {
	u, err := user.Current()
	if err != nil {
		return
	}

	data, err := json.Marshal(c)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf(configFile, u.HomeDir, engineName), data, 0644)
	if err != nil {
		log.Println(err)
	}

}
