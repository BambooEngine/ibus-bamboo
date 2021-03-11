package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
)

import (
	"github.com/dkolbly/wl"
	"github.com/dkolbly/wl/ui"
)

func init() {
	flag.Parse()
	log.SetFlags(0)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	exitChan := make(chan bool, 10)

	if flag.NArg() == 0 {
		log.Fatalf("usage: %s imagefile", os.Args[0])
	}

	img, err := ImageFromFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	display, err := ui.Connect("")
	if err != nil {
		log.Fatal(err)
	}

	b := img.Bounds()
	w := int32(b.Dx())
	h := int32(b.Dy())

	window, err := display.NewWindow(w, h)
	if err != nil {
		log.Fatal(err)
	}

	display.Keyboard().AddKeyHandler(quitter{exitChan})

	window.Draw(img)

loop:
	for {
		select {
		case <-exitChan:
			break loop
		case display.Dispatch() <- struct{}{}:
		}
	}

	log.Print("Loop finished")
	window.Dispose()
	display.Disconnect()
}

type quitter struct {
	ch chan bool
}

func (q quitter) HandleKeyboardKey(ev wl.KeyboardKeyEvent) {
	if ev.Key == 16 {
		q.ch <- true
	}
}
