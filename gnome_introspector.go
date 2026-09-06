package main

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

// Note: dbus.SessionBus is a process-shared connection: it must neither be
// re-helloed nor closed here. The previous code did both on every FocusIn,
// churning (or leaking) a bus connection per focus event on GNOME.
func gnomeGetFocusWindowClass() (string, error) {
	conn, err := dbus.SessionBus()
	var s string
	if err != nil {
		return s, err
	}

	js_code := "global.get_window_actors().find(window => !Main.overview.visible && window.meta_window.has_focus()).get_meta_window().get_wm_class()"
	obj := conn.Object("org.gnome.Shell", "/org/gnome/Shell")
	var ok bool
	err = obj.Call("org.gnome.Shell.Eval", 0, js_code).Store(&ok, &s)
	if !ok {
		if isGnomeOverviewVisible(conn) {
			return "org.gnome.Overview", nil
		} else {
			err = errors.New(s)
		}
	}
	if err != nil {
		return "", err
	}
	return s, nil
}

func isGnomeOverviewVisible(conn *dbus.Conn) bool {
	js_code := "Main.overview.visible"
	obj := conn.Object("org.gnome.Shell", "/org/gnome/Shell")
	var visible string
	var ok bool
	err := obj.Call("org.gnome.Shell.Eval", 0, js_code).Store(&ok, &visible)
	if !ok || err != nil {
		return false
	}
	return visible == "true"
}
