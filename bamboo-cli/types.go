package main

import "flag"

type BambooCliCommand struct {
	name     string
	help     string
	flags    *flag.FlagSet
	dispatch func(args []string) error
}
