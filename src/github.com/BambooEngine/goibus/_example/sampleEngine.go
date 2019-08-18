package main

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/sarim/goibus/ibus"
)

type GittuEngine struct {
	ibus.Engine
	propList *ibus.PropList
}

func (e *GittuEngine) ProcessKeyEvent(keyval uint32, keycode uint32, state uint32) (bool, *dbus.Error) {
	fmt.Println("Process Key Event > ", keyval, keycode, state)

	if state == 0 && keyval == 115 {
		e.UpdateAuxiliaryText(ibus.NewText("s"), true)

		lt := ibus.NewLookupTable()
		lt.AppendCandidate("sss")
		lt.AppendCandidate("s")
		lt.AppendCandidate("gittu")
		lt.AppendLabel("1:")
		lt.AppendLabel("2:")
		lt.AppendLabel("3:")

		e.UpdateLookupTable(lt, true)

		e.UpdatePreeditText(ibus.NewText("s"), uint32(1), true)
		return true, nil
	}
	return false, nil
}

func (e *GittuEngine) FocusIn() *dbus.Error {
	fmt.Println("FocusIn")
	e.RegisterProperties(e.propList)
	return nil
}

func (e *GittuEngine) PropertyActivate(prop_name string, prop_state uint32) *dbus.Error {
	fmt.Println("PropertyActivate", prop_name)
	return nil
}

var eid = 0

func GittuEngineCreator(conn *dbus.Conn, engineName string) dbus.ObjectPath {
	eid++
	fmt.Println("Creating Gittu Engine #", eid)
	objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/IBus/Engine/GittuGo/%d", eid))

	propp := ibus.NewProperty(
		"setup",
		ibus.PROP_TYPE_NORMAL,
		"Preferences - Gittu",
		"gtk-preferences",
		"Configure Gittu Engine",
		true,
		true,
		ibus.PROP_STATE_UNCHECKED)

	engine := &GittuEngine{ibus.BaseEngine(conn, objectPath), ibus.NewPropList(propp)}
	ibus.PublishEngine(conn, objectPath, engine)
	return objectPath
}
