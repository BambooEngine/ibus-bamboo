package main

import (
	"fmt"
	"sort"

	"ibus-bamboo/config"
	"ibus-bamboo/introspectors"
)

var introspectorCommands []BambooCliCommand

var introspectorHelpCommand = BambooCliCommand{
	name: "help",
	help: "Show help for introspector commands",
	dispatch: func(args []string) error {
		printHelp("bamboo-cli introspector - Manage the focused app introspector", introspectorCommands)
		return nil
	},
}

var introspectorListCommand = BambooCliCommand{
	name: "list",
	help: "List available introspectors",
	dispatch: func(args []string) error {
		var names []string
		for _, name := range introspectors.IntrospectorNames {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Println(name)
		}
		return nil
	},
}

var introspectorSetCommand = BambooCliCommand{
	name: "set",
	help: "Set introspector value",
	dispatch: func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: introspector set <introspector>")
		}
		introspector, ok := introspectors.IntrospectorByName[args[0]]
		if !ok {
			return fmt.Errorf("unknown introspector %q", args[0])
		}
		cfg := config.LoadConfig(engineName)
		cfg.Introspector = introspector
		config.SaveConfig(cfg, engineName)
		return nil
	},
}

var introspectorSetCommandCommand = BambooCliCommand{
	name: "set-command",
	help: "Set introspector command",
	dispatch: func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: introspector set-command <command>")
		}
		cfg := config.LoadConfig(engineName)
		cfg.IntrospectorCommand = args[0]
		config.SaveConfig(cfg, engineName)
		return nil
	},
}

var introspectorGetCommand = BambooCliCommand{
	name: "get",
	help: "Print current introspector",
	dispatch: func(args []string) error {
		cfg := config.LoadConfig(engineName)
		fmt.Println(cfg.Introspector)
		return nil
	},
}

var introspectorGetCommandCommand = BambooCliCommand{
	name: "get-command",
	help: "Print current introspector command",
	dispatch: func(args []string) error {
		cfg := config.LoadConfig(engineName)
		fmt.Println(cfg.IntrospectorCommand)
		return nil
	},
}

var introspectorGetFocusedCommand = BambooCliCommand{
	name: "get-focused",
	help: "Get currently focused app id",
	dispatch: func(args []string) error {
		cfg := config.LoadConfig(engineName)
		introspectors.BootstrapIntrospector(cfg)
		fmt.Println(introspectors.GetLatestAppId(cfg.Introspector))
		return nil
	},
}

var introspectorCommand = BambooCliCommand{
	name: "introspector",
	help: "Manage the focused app introspector",
	dispatch: func(args []string) error {
		if len(args) == 0 {
			introspectorHelpCommand.dispatch(args)
			return fmt.Errorf("missing introspector command")
		}

		var action, rest = args[0], args[1:]
		for _, command := range introspectorCommands {
			if command.name == action {
				return command.dispatch(rest)
			}
		}

		introspectorHelpCommand.dispatch(rest)
		return fmt.Errorf("unknown introspector command %s", action)
	},
}

// Assigned in init to avoid an initialization cycle with introspectorHelpCommand
func init() {
	introspectorCommands = []BambooCliCommand{
		introspectorHelpCommand,
		introspectorListCommand,
		introspectorGetCommand,
		introspectorSetCommand,
		introspectorGetCommandCommand,
		introspectorSetCommandCommand,
		introspectorGetFocusedCommand,
	}
}
