/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
	"log"
)

const (
	ComponentName = "org.freedesktop.IBus.bamboo"
	EngineName    = "Bamboo"
	HomePage      = "https://github.com/BambooEngine/ibus-bamboo"
	VnConvertPage = "https://font.ssc.vn"
)

func main() {
	if isIBusDaemonChild() {
		bus := ibus.NewBus()
		bus.RequestName(ComponentName, 0)

		conn := bus.GetDbusConn()
		ibus.NewFactory(conn, IBusBambooEngineCreator)

		select {}
	} else {
		log.Println("Running debug mode")
		runMode = " (debug)"

		bus := ibus.NewBus()
		bus.RegisterComponent(makeDebugComponent())

		conn := bus.GetDbusConn()
		ibus.NewFactory(conn, IBusBambooEngineCreator)

		log.Println("Setting Global Engine to", DebugEngineName)
		bus.CallMethod("SetGlobalEngine", 0, DebugEngineName)

		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)

		select {
		case <-c:
		}
	}
}
