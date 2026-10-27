package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// printHelp pretty prints a heading followed by the commands and their help
func printHelp(heading string, commands []BambooCliCommand) {
	fmt.Println(heading)
	fmt.Println()
	fmt.Println("Commands:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	for _, command := range commands {
		fmt.Fprintf(w, "  %s\t%s\n", command.name, command.help)
	}
	w.Flush()
}
