package ui

import (
	"fmt"
	"image"
	"image/draw"
	//"log"
	"syscall"

	"github.com/dkolbly/wl"
	"github.com/dkolbly/wl/xdg"
)

type Window struct {
	display    *Display
	surface    *wl.Surface
	shSurface  *wl.ShellSurface
	xdgSurface *xdg.Surface
	buffer     *wl.Buffer
	data       []byte
	image      *BGRA
	title      string
	pending    Config
	current    Config
}

func (d *Display) NewWindow(width, height int32) (*Window, error) {
	var err error
	stride := width * 4

	w := new(Window)
	pend := Config{
		Width:  int(width),
		Height: int(height),
	}

	w.pending = pend
	w.current = pend

	w.display = d

	w.surface, err = d.compositor.CreateSurface()
	if err != nil {
		return nil, fmt.Errorf("Surface creation failed: %s", err)
	}

	w.buffer, w.data, err = d.newBuffer(width, height, stride)
	if err != nil {
		return nil, err
	}

	if d.wmBase != nil {
		// New XDG shell
		w.setupXDGTopLevel()
	} else {
		// older plain-jane wl_shell
		w.shSurface, err = d.shell.GetShellSurface(w.surface)
		if err != nil {
			return nil, fmt.Errorf("Shell.GetShellSurface failed: %s", err)
		}

		w.shSurface.AddPingHandler(w)
		w.shSurface.SetToplevel()
	}

	err = w.surface.Attach(w.buffer, width, height)
	if err != nil {
		return nil, fmt.Errorf("Surface.Attach failed: %s", err)
	}

	err = w.surface.Damage(0, 0, width, height)
	if err != nil {
		return nil, fmt.Errorf("Surface.Damage failed: %s", err)
	}

	if true {
		err = w.surface.Commit()
		if err != nil {
			return nil, fmt.Errorf("Surface.Commit failed: %s", err)
		}

		w.image = NewBGRAWithData(
			image.Rect(0, 0, int(width), int(height)),
			w.data)

		d.registerWindow(w)
	}

	return w, nil
}

func (w *Window) DrawUsingFunc(fn func(*BGRA)) {
	fn(w.image)
}

func (w *Window) Draw(img image.Image) {
	draw.Draw(w.image, img.Bounds(), img, img.Bounds().Min, draw.Src)
}

func (w *Window) Dispose() {
	if w.shSurface != nil {
		w.shSurface.RemovePingHandler(w)
	}
	w.surface.Destroy()
	w.buffer.Destroy()
	syscall.Munmap(w.data)
	w.display.unregisterWindow(w)
}

func (w *Window) HandleShellSurfacePing(ev wl.ShellSurfacePingEvent) {
	w.shSurface.Pong(ev.Serial)
}
