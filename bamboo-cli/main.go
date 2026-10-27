package main

import (
	"fmt"
	"os"
)

const engineName = "bamboo"

// Set at build time with -ldflags "-X main.Version=..."
var Version string

var helpCommand = BambooCliCommand{
	name: "help",
	help: "Show this help message",
	dispatch: func(args []string) error {
		printHelp("bamboo-cli - Configuration tool for ibus-bamboo", commands)
		return nil
	},
}

var versionCommand = BambooCliCommand{
	name: "version",
	help: "Show ibus-bamboo version",
	dispatch: func(args []string) error {
		fmt.Printf("Ibus Bamboo %s\n", Version)
		return nil
	},
}

var commands []BambooCliCommand

// Assigned in init to avoid an initialization cycle with helpCommand
func init() {
	commands = []BambooCliCommand{helpCommand, introspectorCommand, modeCommand, versionCommand}
}

func main() {
	var args = os.Args
	if len(args) == 1 {
		helpCommand.dispatch(args)
		os.Exit(2)
	}

	var action, rest = os.Args[1], os.Args[2:]
	for _, command := range commands {
		if command.name == action {
			if err := command.dispatch(rest); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	fmt.Printf("Unknown action %s\n\n", action)
	helpCommand.dispatch(rest)
	os.Exit(1)
}
