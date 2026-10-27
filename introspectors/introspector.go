package introspectors

import (
	"os"
	"strings"

	"ibus-bamboo/types"
)

const (
	// The CurrentClientId state is stored in this module
	// A subprocess will be spawn to watch over the currently
	// active client
	Reactive types.IntrospectorMode = iota

	// The current client id will be retreived by IBus Bamboo
	// when it is needed, by some kind of API. No subprocess is needed for this one.
	Passive
)

const (
	IntrospectorAuto    types.BambooIntrospector = "auto"
	IntrospectorX11     types.BambooIntrospector = "x11"
	IntrospectorGnome   types.BambooIntrospector = "gnome"
	IntrospectorWayland types.BambooIntrospector = "wayland"
	IntrospectorCommand types.BambooIntrospector = "command"
)

var IntrospectorNames = map[types.BambooIntrospector]string{
	IntrospectorAuto:    "auto",
	IntrospectorX11:     "x11",
	IntrospectorGnome:   "gnome",
	IntrospectorWayland: "wayland",
	IntrospectorCommand: "command",
}

var IntrospectorByName = map[string]types.BambooIntrospector{
	"":        IntrospectorAuto,
	"auto":    IntrospectorAuto,
	"x11":     IntrospectorX11,
	"gnome":   IntrospectorGnome,
	"wayland": IntrospectorWayland,
	"command": IntrospectorCommand,
}

// GetIntrospectorMode tells whether an introspector keeps the focused window
// class up to date in the background (Reactive) or queries it on demand (Passive).
func GetIntrospectorMode(introspector types.BambooIntrospector) types.IntrospectorMode {
	switch introspector {
	case IntrospectorWayland:
		// WlGetFocusWindowClass runs in a goroutine and keeps WlAppId updated
		return Reactive
	case IntrospectorX11, IntrospectorGnome, IntrospectorCommand:
		return Passive
	default:
		// IntrospectorAuto currently falls back to querying X11 on demand
		return Passive
	}
}

func hasGnome(env string) bool {
	return strings.Contains(strings.ToLower(os.Getenv(env)), "gnome")
}

// guessIntrospector picks an introspector from the running session.
// GNOME is checked first: GNOME on Wayland does not implement the
// wlr foreign toplevel protocol used by the Wayland introspector.
func guessIntrospector() types.BambooIntrospector {
	isWayland := os.Getenv("WAYLAND_DISPLAY") != ""
	isGnome := hasGnome("XDG_CURRENT_DESKTOP") || hasGnome("DESKTOP_SESSION") || hasGnome("GDMSESSION")
	if isGnome {
		return IntrospectorGnome
	}
	if isWayland {
		return IntrospectorWayland
	}
	return IntrospectorX11
}

// BootstrapIntrospector starts whatever background work an introspector needs
// before the focused window class can be read from it.
// IntrospectorAuto in cfg is replaced by the introspector guessed from the session.
func BootstrapIntrospector(cfg *types.Config) {
	if cfg.Introspector == IntrospectorAuto {
		cfg.Introspector = guessIntrospector()
	}
	switch cfg.Introspector {
	case IntrospectorWayland:
		go WlGetFocusWindowClass()
	case IntrospectorCommand:
		CmdSetIntrospectCommand(cfg.IntrospectorCommand)
	}
}

// GetLatestAppId returns the class of the currently focused window,
// falling back to X11 when the introspector cannot tell.
func GetLatestAppId(introspector types.BambooIntrospector) string {
	var wmClass string
	switch introspector {
	case IntrospectorGnome:
		wmClass, _ = GnomeGetFocusWindowClass()
	case IntrospectorWayland:
		wmClass = WlAppId
	case IntrospectorX11:
		wmClass = X11GetFocusWindowClass()
	case IntrospectorCommand:
		wmClass = CmdGetFocusedWmClass()
	}

	if wmClass == "" {
		wmClass = X11GetFocusWindowClass()
	}
	wmClass = strings.Replace(wmClass, "\"", "", -1)
	return wmClass
}
