package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"ibus-bamboo/config"
	"ibus-bamboo/introspectors"
)

var modeCommands []BambooCliCommand

// parseMode accepts either the integer value or the short name of a typing mode
func parseMode(arg string) (int, error) {
	if mode, err := strconv.Atoi(arg); err == nil {
		if _, ok := config.ImLookupTable[mode]; ok {
			return mode, nil
		}
	}
	for mode, name := range config.ImShortNames {
		if name == arg {
			return mode, nil
		}
	}
	return 0, fmt.Errorf("unknown typing mode %q", arg)
}

var modeHelpCommand = BambooCliCommand{
	name: "help",
	help: "Show help for typing mode commands",
	dispatch: func(args []string) error {
		printHelp("bamboo-cli mode - Manage the typing mode", modeCommands)
		return nil
	},
}

var modeListCommand = BambooCliCommand{
	name: "list",
	help: "List all typing modes (--json for JSON output)",
	dispatch: func(args []string) error {
		flags := flag.NewFlagSet("mode list", flag.ContinueOnError)
		asJSON := flags.Bool("json", false, "Output as JSON")
		if err := flags.Parse(args); err != nil {
			return err
		}
		var modes []int
		for mode := range config.ImLookupTable {
			modes = append(modes, mode)
		}
		sort.Ints(modes)
		if *asJSON {
			type modeEntry struct {
				Index       int    `json:"index"`
				Name        string `json:"name"`
				DisplayName string `json:"display_name"`
			}
			entries := []modeEntry{}
			for _, mode := range modes {
				entries = append(entries, modeEntry{mode, config.ImShortNames[mode], config.ImLookupTable[mode]})
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			return enc.Encode(entries)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		for _, mode := range modes {
			fmt.Fprintf(w, "%d\t%s\t%s\n", mode, config.ImShortNames[mode], config.ImLookupTable[mode])
		}
		w.Flush()
		return nil
	},
}

var modeGetCommand = BambooCliCommand{
	name: "get",
	help: "Get current typing mode (for focused client)",
	dispatch: func(args []string) error {
		cfg := config.LoadConfig(engineName)
		introspectors.BootstrapIntrospector(cfg)
		wmClass := introspectors.GetLatestAppId(cfg.Introspector)

		// Same resolution as the engine's getInputMode
		mode := config.PreeditIM
		if im, ok := cfg.InputModeMapping[wmClass]; ok && wmClass != "" && config.ImLookupTable[im] != "" {
			mode = im
		} else if _, ok := config.ImLookupTable[cfg.DefaultInputMode]; ok {
			mode = cfg.DefaultInputMode
		}
		fmt.Println(config.ImShortNames[mode])
		return nil
	},
}

var modeGetDefaultCommand = BambooCliCommand{
	name: "get-default",
	help: "Print current default typing mode",
	dispatch: func(args []string) error {
		cfg := config.LoadConfig(engineName)
		fmt.Println(config.ImShortNames[cfg.DefaultInputMode])
		return nil
	},
}

var modeSetDefaultCommand = BambooCliCommand{
	name: "set-default",
	help: "Set default typing mode (integer or short name)",
	dispatch: func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: mode set-default <mode>")
		}
		mode, err := parseMode(args[0])
		if err != nil {
			return err
		}
		cfg := config.LoadConfig(engineName)
		cfg.DefaultInputMode = mode
		config.SaveConfig(cfg, engineName)
		fmt.Println(config.ImShortNames[mode])
		return nil
	},
}

func modeSetHelp() error {
	return fmt.Errorf(`usage:
	mode set <app_id/wm_class> <mode>
	mode set --focused <mode>`)
}

var modeSetCommand = BambooCliCommand{
	name: "set",
	help: "Set typing mode for an app id/wm_class, or the focused app with --focused (integer or short name)",
	dispatch: func(args []string) error {
		flags := flag.NewFlagSet("mode set", flag.ContinueOnError)
		focused := flags.Bool("focused", false, "Use the app id/wm_class of the focused client")
		if err := flags.Parse(args); err != nil {
			return err
		}
		args = flags.Args()
		cfg := config.LoadConfig(engineName)
		var appId string
		if *focused {
			if len(args) != 1 {
				return modeSetHelp()
			}
			introspectors.BootstrapIntrospector(cfg)
			appId = introspectors.GetLatestAppId(cfg.Introspector)
			if appId == "" {
				return fmt.Errorf("could not determine the focused app id/wm_class")
			}
		} else {
			if len(args) != 2 {
				return modeSetHelp()
			}
			appId = args[0]
		}
		mode, err := parseMode(args[len(args)-1])
		if err != nil {
			return err
		}
		cfg.InputModeMapping[appId] = mode
		config.SaveConfig(cfg, engineName)
		fmt.Println(config.ImShortNames[mode])
		return nil
	},
}

var modeCommand = BambooCliCommand{
	name: "mode",
	help: "Manage the typing mode",
	dispatch: func(args []string) error {
		if len(args) == 0 {
			modeHelpCommand.dispatch(args)
			return fmt.Errorf("missing mode command")
		}

		var action, rest = args[0], args[1:]
		for _, command := range modeCommands {
			if command.name == action {
				return command.dispatch(rest)
			}
		}

		modeHelpCommand.dispatch(rest)
		return fmt.Errorf("unknown mode command %s", action)
	},
}

// Assigned in init to avoid an initialization cycle with modeHelpCommand
func init() {
	modeCommands = []BambooCliCommand{
		modeHelpCommand,
		modeListCommand,
		modeGetCommand,
		modeSetCommand,
		modeGetDefaultCommand,
		modeSetDefaultCommand,
	}
}
