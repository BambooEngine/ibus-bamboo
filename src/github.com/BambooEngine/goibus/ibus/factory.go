package ibus

import (
	"github.com/godbus/dbus"
)

type Factory struct {
	conn          *dbus.Conn
	EngineCreator func(conn *dbus.Conn, engineName string) dbus.ObjectPath
}

func NewFactory(conn *dbus.Conn, EngineCreator func(conn *dbus.Conn, engineName string) dbus.ObjectPath) *Factory {
	f := &Factory{conn, EngineCreator}
	conn.Export(f, "/org/freedesktop/IBus/Factory", IBUS_IFACE_ENGINE_FACTORY)
	return f
}

// # Return a array. [name, default_language, icon_path, authors, credits]
// @method(out_signature="as")
// def GetInfo(self): pass

// # Factory should allocate all resources in this method
// @method()
// def Initialize(self): pass

// # Factory should free all allocated resources in this method
// @method()
// def Uninitialize(self): pass

// # Create an input context and return the id of the context.
// # If failed, it will return "" or None.
// @method(in_signature="s", out_signature="o")
func (factory *Factory) CreateEngine(engineName string) (dbus.ObjectPath, *dbus.Error) {
	return factory.EngineCreator(factory.conn, engineName), nil
}

// # Destroy the engine
// @method()
func (factory *Factory) Destroy() *dbus.Error {
	factory.conn.Export(nil, "/org/freedesktop/IBus/Factory", IBUS_IFACE_ENGINE_FACTORY)
	return nil
}
