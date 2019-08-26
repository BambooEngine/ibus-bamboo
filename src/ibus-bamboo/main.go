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
	"flag"
	"fmt"
	"github.com/BambooEngine/bamboo-core"
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
	"log"
	"os"
)

const (
	ComponentName = "org.freedesktop.IBus.bamboo"
	EngineName    = "Bamboo"
)

var embedded = flag.Bool("ibus", false, "Run the embedded ibus component")
var version = flag.Bool("version", false, "Show version")

func init() {
	flag.Parse()
	if *embedded {
		os.Chdir(DataDir)
	}
	go func() {
		loadEmojiOne(DictEmojiOne)
		var dictionary, _ = loadDictionary(DictVietnameseCm)
		bamboo.AddDictionaryToSpellingTrie(dictionary)
	}()
}

func main() {
	if *version {
		fmt.Println(Version)
	} else if *embedded {
		engine := GetBambooEngineCreator()
		bus := ibus.NewBus()
		bus.RequestName(ComponentName, 0)

		conn := bus.GetDbusConn()
		ibus.NewFactory(conn, engine)

		select {}
	} else {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
		bus := ibus.NewBus()
		log.Println("Got Bus, Running Standalone")
		component := &ibus.Component{
			Name:          "IBusComponent",
			ComponentName: ComponentName + "Standalone",
		}
		engine := &ibus.EngineDesc{
			Name:       "IBusEngineDesc",
			EngineName: EngineName + "Standalone",
		}
		component.AddEngine(engine)
		bus.RegisterComponent(component)

		conn := bus.GetDbusConn()
		ibus.NewFactory(conn, GetBambooEngineCreator())

		bus.CallMethod("SetGlobalEngine", 0, EngineName+"Standalone")

		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)

		select {
		case <-c:
		}
	}
}
