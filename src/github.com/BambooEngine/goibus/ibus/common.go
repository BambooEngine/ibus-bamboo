package ibus

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
)

const (
	BUS_DAEMON_NAME     = "org.freedesktop.DBus"
	BUS_DAEMON_PATH     = "/org/freedesktop/DBus"
	BUS_PROPERTIES_NAME = "org.freedesktop.DBus.Properties"

	IBUS_IFACE_IBUS   = "org.freedesktop.IBus"
	IBUS_PATH_IBUS    = "/org/freedesktop/IBus"
	IBUS_SERVICE_IBUS = "org.freedesktop.IBus"

	IBUS_IFACE_PANEL          = "org.freedesktop.IBus.Panel"
	IBUS_IFACE_CONFIG         = "org.freedesktop.IBus.Config"
	IBUS_IFACE_SERVICE        = "org.freedesktop.IBus.Service"
	IBUS_IFACE_ENGINE         = "org.freedesktop.IBus.Engine"
	IBUS_IFACE_ENGINE_FACTORY = "org.freedesktop.IBus.Factory"
	IBUS_IFACE_INPUT_CONTEXT  = "org.freedesktop.IBus.InputContext"
	IBUS_IFACE_NOTIFICATIONS  = "org.freedesktop.IBus.Notifications"

	IBUS_ENGINE_PREEDIT_CLEAR  uint32 = 0
	IBUS_ENGINE_PREEDIT_COMMIT uint32 = 1

	ORIENTATION_HORIZONTAL int32 = 0
	ORIENTATION_VERTICAL   int32 = 1
	ORIENTATION_SYSTEM     int32 = 2

	PROP_TYPE_NORMAL    uint32 = 0
	PROP_TYPE_TOGGLE    uint32 = 1
	PROP_TYPE_RADIO     uint32 = 2
	PROP_TYPE_MENU      uint32 = 3
	PROP_TYPE_SEPARATOR uint32 = 4

	PROP_STATE_UNCHECKED    uint32 = 0
	PROP_STATE_CHECKED      uint32 = 1
	PROP_STATE_INCONSISTENT uint32 = 2

	IBUS_ATTR_TYPE_NONE       uint32 = 0
	IBUS_ATTR_TYPE_UNDERLINE  uint32 = 1
	IBUS_ATTR_TYPE_FOREGROUND uint32 = 2
	IBUS_ATTR_TYPE_BACKGROUND uint32 = 3

	IBUS_ATTR_UNDERLINE_NONE   uint32 = 0
	IBUS_ATTR_UNDERLINE_SINGLE uint32 = 1
	IBUS_ATTR_UNDERLINE_DOUBLE uint32 = 2
	IBUS_ATTR_UNDERLINE_LOW    uint32 = 3
	IBUS_ATTR_UNDERLINE_ERROR  uint32 = 4
)

func GetAddress() string {
	address := os.Getenv("IBUS_ADDRESS")
	if address != "" {
		return address
	}
	data, err := ioutil.ReadFile(GetSocketPath())
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.Index(line, "IBUS_ADDRESS=") == 0 {
			address = line[13:]
		}
	}
	return address
}

func GetSocketPath() string {
	path := os.Getenv("IBUS_ADDRESS_FILE")
	if path != "" {
		return path
	}
	display := os.Getenv("DISPLAY")
	if display == "" {
		fmt.Fprintf(os.Stderr, "DISPLAY is empty! We use default DISPLAY (:0.0)")
		display = ":0.0"
	}
	// format is {hostname}:{displaynumber}.{screennumber}
	hostname := "unix"
	HDS := strings.SplitN(display, ":", 2)
	DS := strings.SplitN(HDS[1], ".", 2)

	if HDS[0] != "" {
		hostname = HDS[0]
	}
	p := fmt.Sprintf("%s-%s-%s", GetLocalMachineId(), hostname, DS[0])
	path = GetUserConfigDir() + "/ibus/bus/" + p

	return path
}

func GetLocalMachineId() string {
	var mID []byte
	var err error
	mID, err = ioutil.ReadFile("/var/lib/dbus/machine-id")
	if err != nil {
		mID, err = ioutil.ReadFile("/etc/machine-id")
		if err != nil {
			panic(err)
		}
	}
	return strings.TrimSpace(string(mID))
}

func GetUserConfigDir() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		return os.Getenv("HOME") + "/.config"
	}
	return dir
}

func GetUserAuth() []dbus.Auth {
	uid := os.Getenv("DBUS_AUTH_UID")
	if uid == "" {
		uid = strconv.Itoa(os.Getuid())
	}
	home := os.Getenv("DBUS_AUTH_HOME")
	if home == "" {
		home = os.Getenv("HOME")
	}
	return []dbus.Auth{dbus.AuthExternal(uid), dbus.AuthCookieSha1(uid, home)}
}
