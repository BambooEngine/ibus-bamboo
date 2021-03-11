package ui

import (
	"fmt"
	//"github.com/dkolbly/wl"
	"github.com/dkolbly/wl/xdg"
)

type Config struct {
	Width  int
	Height int
	Active bool
}

// see description of xdg_surface.configure
//
// Basically, we will receive a series of events on our role
// object (e.g., a xdg.Toplevel) which are accumulating latchable
// state.  When the xdg.Surface get a configure event, we "latch"
// those changes, do whatever we need to do, and then respond with
// an AckConfigure request.

// the compositor wants the surface to be closed, based on user action
func (w *Window) HandleToplevelClose(ev xdg.ToplevelCloseEvent) {
	fmt.Printf("toplevel close request: %v\n", ev)
}

func (w *Window) HandleToplevelConfigure(ev xdg.ToplevelConfigureEvent) {
	fmt.Printf("toplevel configured: %v\n", ev)

	pend := Config{
		Width:  int(ev.Width),
		Height: int(ev.Height),
	}
	for _, state := range ev.States {
		if state == xdg.ToplevelStateActivated {
			pend.Active = true
		}
	}
	w.pending = pend
}

func (w *Window) HandleSurfaceConfigure(ev xdg.SurfaceConfigureEvent) {
	fmt.Printf("surface configured; committing %#v\n", w.pending)
	// and ack it
	w.xdgSurface.AckConfigure(ev.Serial)

	// apply the changes
	w.current = w.pending
}

func (w *Window) setupXDGTopLevel() error {

	d := w.display

	/*fmt.Printf("creating xdg_wm_base\n")
	wm := xdg.NewXdgWmBase(d.Context())
	fmt.Printf("==> %#v\n", wm)
	*/
	s, err := d.wmBase.GetXdgSurface(w.surface)
	if err != nil {
		fmt.Printf("failed to get surface: %s", err)
	} else {
		fmt.Printf("surface is: %p\n", s)
	}

	w.xdgSurface = s
	/*ping := wl.HandlerFunc(func(x interface{}) {
		if ev, ok := x.(xdg.WmBasePingEvent); ok {
			fmt.Printf("ping <%d>\n", ev.Serial)
			d.wmBase.Pong(ev.Serial)
		} else {
			fmt.Printf("umm, what %#v\n", x)
		}
	})
	foo := wl.HandlerFunc(func(x interface{}) {
		fmt.Printf("surface configured: %#v\n", x)
	})*/

	s.AddConfigureHandler(w)

	top, err := s.GetToplevel()
	if err != nil {
		panic(err)
	}
	fmt.Printf("top level is: %p\n", top)

	top.SetTitle("Hello!")
	top.SetAppId("go.hello")
	/*bar := wl.HandlerFunc(func(x interface{}) {
		fmt.Printf("toplevel configured: %#v\n", x)
	})
	barc := wl.HandlerFunc(func(x interface{}) {
		fmt.Printf("toplevel closed: %#v\n", x)
	})*/
	top.AddConfigureHandler(w)
	top.AddCloseHandler(w)

	err = s.SetWindowGeometry(10, 10, 300, 300)
	if err != nil {
		panic(err)
	}
	// we need to commit the underlying wl_surface before
	// doing much else (see description of xdg_surface)
	err = w.surface.Commit()
	if err != nil {
		return fmt.Errorf("Surface.Commit failed: %s", err)
	}

	fmt.Printf("ok...\n")
	return nil
}
