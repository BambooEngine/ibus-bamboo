package main

import (
	"github.com/godbus/dbus"
)

func gnomeGetFocusWindowClass() string {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	js_code := `
        global._ib_current_window = () => {
            var window_list = global.get_window_actors();
            var active_window_actor = window_list.find(window => window.meta_window.has_focus());
            var active_window = active_window_actor.get_meta_window();
            var vm_class = active_window.get_wm_class();
            var title = active_window.get_title();
            var result = vm_class;
            return result;
        }
				`
	obj := conn.Object("org.gnome.Shell", "/org/gnome/Shell")
	call := obj.Call("org.gnome.Shell.Eval", 0, js_code)
	if call.Err != nil {
		panic(call.Err)
	}
	var s string
	var ok bool
	err = obj.Call("org.gnome.Shell.Eval", 0, "global._ib_current_window()").Store(&ok, &s)
	if (err != nil) {
		panic(err)
	}
	return s
}
