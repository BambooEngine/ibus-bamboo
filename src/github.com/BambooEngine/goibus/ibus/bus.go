package ibus

import (
	"github.com/godbus/dbus"
)

type Bus struct {
	dbusConn   *dbus.Conn
	dbusObject dbus.BusObject
	ibusObject dbus.BusObject
}

func NewBus() *Bus {
	doPanic := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	addr := GetAddress()
	conn, err := dbus.Dial(addr)
	doPanic(err)

	err = conn.Auth(GetUserAuth())
	doPanic(err)

	err = conn.Hello()
	doPanic(err)

	dbusObject := conn.Object(BUS_DAEMON_NAME, dbus.ObjectPath(BUS_DAEMON_PATH))
	ibusObject := conn.Object(IBUS_SERVICE_IBUS, dbus.ObjectPath(IBUS_PATH_IBUS))

	return &Bus{conn, dbusObject, ibusObject}
}

func (bus *Bus) CallMethod(name string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	return bus.ibusObject.Call(bus.ibusObject.Destination()+"."+name, flags, args...)
}

func (bus *Bus) RequestName(name string, flags dbus.RequestNameFlags) (dbus.RequestNameReply, error) {
	return bus.dbusConn.RequestName(name, flags)
}

func (bus *Bus) RegisterComponent(component *Component) {
	bus.CallMethod("RegisterComponent", 0, dbus.MakeVariant(component))
}

func (bus *Bus) GetDbusConn() *dbus.Conn {
	return bus.dbusConn
}
