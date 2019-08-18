package main

import (
	"flag"
	"fmt"
	"github.com/godbus/dbus"
	"github.com/sarim/goibus/ibus"
	"os"
)

var embeded = flag.Bool("ibus", false, "Run the embeded ibus component")
var standalone = flag.Bool("standalone", false, "Run standalone by creating new component")
var generatexml = flag.String("xml", "", "Write xml representation of component to file or stdout if file == \"-\"")

func makeComponent() *ibus.Component {

	component := ibus.NewComponent(
		"org.freedesktop.IBus.bamboo",
		"Gittu Sample",
		"2.0",
		"MPL 1.1",
		"Sarim Khan <sarim2005@gmail.com>",
		"https://github.com/sarim/goibus",
		"/usr/bin/gittuengine",
		"gittu-sample")

	avroenginedesc := ibus.SmallEngineDesc(
		"gittu-sample",
		"Gittu Sample",
		"Gittu Sample Engine",
		"en",
		"MPL 1.1",
		"Sarim Khan <sarim2005@gmail.com>",
		"/usr/share/gittu/icon.png",
		"en",
		"/usr/bin/gittupref",
		"2.0")

	component.AddEngine(avroenginedesc)

	return component
}

func main() {

	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.CommandLine.VisitAll(func(f *flag.Flag) {
			format := "  -%s: %s\n"
			fmt.Fprintf(os.Stderr, format, f.Name, f.Usage)
		})
	}

	flag.Parse()

	if *generatexml != "" {
		c := makeComponent()

		if *generatexml == "-" {
			c.OutputXML(os.Stdout)
		} else {
			f, err := os.Create(*generatexml)
			if err != nil {
				panic(err)
			}

			c.OutputXML(f)
			f.Close()
		}
	} else if *embeded {
		bus := ibus.NewBus()
		fmt.Println("Got Bus, Running Embeded")

		conn := bus.GetDbusConn()
		ibus.NewFactory(conn, GittuEngineCreator)
		bus.RequestName("org.freedesktop.IBus.bamboo", 0)
		select {}
	} else if *standalone {
		bus := ibus.NewBus()
		fmt.Println("Got Bus, Running Standalone")

		conn := bus.GetDbusConn()
		ibus.NewFactory(conn, GittuEngineCreator)
		bus.RegisterComponent(makeComponent())

		fmt.Println("Setting Global Engine to me")
		bus.CallMethod("SetGlobalEngine", 0, "gittu-sample")

		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)

		select {
		case <-c:
		}

	} else {
		Usage()
		os.Exit(1)
	}
}
