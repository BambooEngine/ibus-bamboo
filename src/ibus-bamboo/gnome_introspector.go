package main

import (
	"errors"

	"github.com/godbus/dbus"
)

func gnomeGetFocusWindowClass() ( string, error ) {
	conn, err := dbus.SessionBus()
	var s string
	if err != nil {
		return s, err
	}
	defer func() {
		if err = conn.Hello(); err == nil {
			conn.Close()
		}
	}()

	js_code := "global.get_window_actors().find(window => window.meta_window.has_focus()).get_meta_window().get_wm_class()"
	obj := conn.Object("org.gnome.Shell", "/org/gnome/Shell")
	var ok bool
	err = obj.Call("org.gnome.Shell.Eval", 0, js_code).Store(&ok, &s)
	if !ok {
		err = errors.New(s)
	}
	if (err != nil) {
		return "", err
	}
	return s, nil
}

func isGnomeOverviewVisible() ( bool ) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false
	}
	defer func() {
		if err = conn.Hello(); err == nil {
			conn.Close()
		}
	}()

	js_code := "Main.overview.visible"
	obj := conn.Object("org.gnome.Shell", "/org/gnome/Shell")
	var visible string
	var ok bool
	err = obj.Call("org.gnome.Shell.Eval", 0, js_code).Store(&ok, &visible)
	if !ok || err != nil {
		return false
	}
	return visible == "true"
}
