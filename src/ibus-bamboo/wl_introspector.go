package main

import (
	"fmt"

	wl "github.com/dkolbly/wl"
)

var wlAppId string

func wlGetFocusWindowClass() error {
	display, err := wl.Connect("")
	if err != nil {
		return fmt.Errorf("Connect to Wayland server failed %s", err)
	}
	appIdChan := make(chan string, 10)
	err = registerGlobals(display, appIdChan)
	if err != nil {
		return err
	}
	for {
		select {
		case wlAppId = <-appIdChan:
			fmt.Println("wlAppId = ", wlAppId)
		case display.Context().Dispatch() <- struct{}{}:
		}
	}
	display.Context().Close()
	return nil
}

func registerGlobals(display *wl.Display, appIdChan chan string) error {
	registry, err := display.GetRegistry()
	if err != nil {
		return fmt.Errorf("Display.GetRegistry failed : %s", err)
	}

	callback, err := display.Sync()
	if err != nil {
		return fmt.Errorf("Display.Sync failed %s", err)
	}

	rgeChan := make(chan wl.RegistryGlobalEvent)
	rgeHandler := registrar{rgeChan}
	registry.AddGlobalHandler(rgeHandler)

	cdeChan := make(chan wl.CallbackDoneEvent)
	cdeHandler := doner{cdeChan}

	callback.AddDoneHandler(cdeHandler)
	adHandler := appIdHandler{appIdChan, false}
loop:
	for {
		select {
		case ev := <-rgeChan:
			if err := registerInterface(registry, ev, display.Context(), &adHandler); err != nil {
				return err
			}
		case display.Context().Dispatch() <- struct{}{}:
		case <-cdeChan:
			break loop
		}
	}

	registry.RemoveGlobalHandler(rgeHandler)
	callback.RemoveDoneHandler(cdeHandler)
	return nil
}

func registerInterface(registry *wl.Registry, ev wl.RegistryGlobalEvent, ctx *wl.Context, adHandler *appIdHandler) error {
	switch ev.Interface {
	case "zwlr_foreign_toplevel_manager_v1":
		manager := NewZwlrForeignToplevelManagerV1(ctx)
		manager.AddToplevelHandler(toplevelHandlers{adHandler})
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, manager)
		if err != nil {
			return fmt.Errorf("Unable to bind ZwlrForeignToplevelManagerV1 interface: %s", err)
		}
	}
	return nil
}

type doner struct {
	ch chan wl.CallbackDoneEvent
}

func (d doner) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	d.ch <- ev
}

type registrar struct {
	ch chan wl.RegistryGlobalEvent
}

func (r registrar) HandleRegistryGlobal(ev wl.RegistryGlobalEvent) {
	r.ch <- ev
}
type toplevelHandlers struct {
	adHandler *appIdHandler
}

func (t toplevelHandlers) HandleZwlrForeignToplevelManagerV1Toplevel(ev ZwlrForeignToplevelManagerV1ToplevelEvent) {
	ev.Toplevel.AddAppIdHandler(t.adHandler)
	ev.Toplevel.AddDoneHandler(t.adHandler)
}

type appIdHandler struct {
	ch chan string
	isAppIdUpdated bool
}

func (a *appIdHandler) HandleZwlrForeignToplevelHandleV1AppId(ev ZwlrForeignToplevelHandleV1AppIdEvent) {
	a.ch <- ev.AppId
	a.isAppIdUpdated = true
}

func (a *appIdHandler) HandleZwlrForeignToplevelHandleV1Done(ev ZwlrForeignToplevelHandleV1DoneEvent) {
	if !a.isAppIdUpdated {
		a.ch <- ""
	}
	a.isAppIdUpdated = false
}

