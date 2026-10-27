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
	"log"
	"os"
	"strings"

	"github.com/BambooEngine/bamboo-core"
	ibus "github.com/BambooEngine/goibus"
	"ibus-bamboo/config"
	"ibus-bamboo/introspectors"
	"ibus-bamboo/types"
)

const (
	ComponentName = "org.freedesktop.IBus.bamboo"
	EngineName    = "Bamboo"
)

var embedded = flag.Bool("ibus", false, "Run the embedded ibus component")
var version = flag.Bool("version", false, "Show version")
var gui = flag.Bool("gui", false, "Show GUI")
var engineRef *IBusBambooEngine
var introspector types.BambooIntrospector

func validateConfig(cfg *types.Config) error {
	resolved, ok := introspectors.IntrospectorByName[string(cfg.Introspector)]
	if !ok {
		return fmt.Errorf("unknown introspector %q", cfg.Introspector)
	}
	cfg.Introspector = resolved
	if cfg.Introspector == introspectors.IntrospectorCommand && cfg.IntrospectorCommand == "" {
		return fmt.Errorf("introspector %q requires IntrospectorCommand", cfg.Introspector)
	}
	if _, ok := cfg.InputMethodDefinitions[cfg.InputMethod]; !ok {
		return fmt.Errorf("unknown input method %q", cfg.InputMethod)
	}
	for _, cs := range bamboo.GetCharsetNames() {
		if cs == cfg.OutputCharset {
			return nil
		}
	}
	return fmt.Errorf("unknown output charset %q", cfg.OutputCharset)
}

func main() {
	flag.Parse()

	/// Print the version and exit
	if *version {
		fmt.Println(Version)
		return
	}

	if *embedded {
		os.Chdir(DataDir)
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.LoadConfig(strings.ToLower(EngineName))

	/// Show configuration that are loaded for debugging
	fmt.Print("Configuration loaded\n")
	fmt.Printf("Config: %s\n", config.GetConfigPath(strings.ToLower(EngineName)))
	fmt.Printf("InputMethod: %s\n", cfg.InputMethod)
	fmt.Printf("InputMethodDefinitions: %d definitions\n", len(cfg.InputMethodDefinitions))
	fmt.Printf("OutputCharset: %s\n", cfg.OutputCharset)
	fmt.Printf("Flags: %d\n", cfg.Flags)
	fmt.Printf("IBflags: %d\n", cfg.IBflags)
	fmt.Printf("Shortcuts: %v\n", cfg.Shortcuts)
	fmt.Printf("DefaultInputMode: %d\n", cfg.DefaultInputMode)
	fmt.Printf("InputModeMapping: %v\n", cfg.InputModeMapping)
	fmt.Printf("Introspector: %v\n", cfg.Introspector)
	fmt.Printf("IntrospectorCommand: %q\n", cfg.IntrospectorCommand)

	if err := validateConfig(cfg); err != nil {
		log.Fatalf("Invalid config %s: %v", config.GetConfigPath(strings.ToLower(EngineName)), err)
	}
	introspectors.BootstrapIntrospector(cfg)
	introspector = cfg.Introspector
	log.Printf("Using introspector: %s", introspector)

	if *embedded {
		log.Println("IbusBamboo running in embedded mode")
		engine := GetIBusEngineCreator()
		bus := ibus.NewBus()
		bus.RequestName(ComponentName, 0)

		conn := bus.GetDbusConn()
		ibus.NewFactory(conn, engine)

		select {}
	} else {
		log.Println("IbusBamboo running in IBus mode")
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
		ibus.NewFactory(conn, GetIBusEngineCreator())

		bus.CallMethod("SetGlobalEngine", 0, EngineName+"Standalone")

		select {}
	}
}
