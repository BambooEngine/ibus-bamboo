package ui

import (
	"fmt"
	"log"
	"sync"
	"syscall"
)

import (
	"github.com/dkolbly/wl"
	"github.com/dkolbly/wl/xdg"
)

type Display struct {
	mu                sync.RWMutex
	display           *wl.Display
	registry          *wl.Registry
	compositor        *wl.Compositor
	subCompositor     *wl.Subcompositor
	shell             *wl.Shell
	shm               *wl.Shm
	seat              *wl.Seat
	dataDeviceManager *wl.DataDeviceManager
	pointer           *wl.Pointer
	keyboard          *wl.Keyboard
	touch             *wl.Touch
	wmBase            *xdg.WmBase
	windows           []*Window
}

func Connect(addr string) (*Display, error) {
	d := new(Display)
	display, err := wl.Connect(addr)
	if err != nil {
		return nil, fmt.Errorf("Connect to Wayland server failed %s", err)
	}

	d.display = display
	display.AddErrorHandler(d)

	err = d.registerGlobals()
	if err != nil {
		return nil, err
	}

	err = d.checkGlobalsRegistered()
	if err != nil {
		return nil, err
	}

	err = d.registerInputs()
	if err != nil {
		return nil, err
	}

	err = d.checkInputsRegistered()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Display) Disconnect() {
	d.keyboard.Release()
	d.pointer.Release()
	if d.touch != nil {
		d.touch.Release()
	}

	d.seat.Release()
	d.display.Context().Close()
}

func (d *Display) Dispatch() chan<- struct{} {
	return d.Context().Dispatch()
}

func (d *Display) Context() *wl.Context {
	return d.display.Context()
}

type registrar struct {
	ch chan wl.RegistryGlobalEvent
}

func (r registrar) HandleRegistryGlobal(ev wl.RegistryGlobalEvent) {
	r.ch <- ev
}

type doner struct {
	ch chan wl.CallbackDoneEvent
}

func (d doner) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	d.ch <- ev
}

func (d *Display) registerGlobals() error {
	registry, err := d.display.GetRegistry()
	if err != nil {
		return fmt.Errorf("Display.GetRegistry failed : %s", err)
	}
	d.registry = registry

	callback, err := d.display.Sync()
	if err != nil {
		return fmt.Errorf("Display.Sync failed %s", err)
	}

	rgeChan := make(chan wl.RegistryGlobalEvent)
	rgeHandler := registrar{rgeChan}
	registry.AddGlobalHandler(rgeHandler)

	cdeChan := make(chan wl.CallbackDoneEvent)
	cdeHandler := doner{cdeChan}

	callback.AddDoneHandler(cdeHandler)
loop:
	for {
		select {
		case ev := <-rgeChan:
			if err := d.registerInterface(registry, ev); err != nil {
				return err
			}
		case d.Dispatch() <- struct{}{}:
		case <-cdeChan:
			break loop
		}
	}

	registry.RemoveGlobalHandler(rgeHandler)
	callback.RemoveDoneHandler(cdeHandler)
	return nil
}

type seatcap struct {
	ch chan wl.SeatCapabilitiesEvent
}

func (sce seatcap) HandleSeatCapabilities(ev wl.SeatCapabilitiesEvent) {
	sce.ch <- ev
}

func (d *Display) registerInputs() error {
	callback, err := d.display.Sync()
	if err != nil {
		return fmt.Errorf("Display.Sync failed %s", err)
	}

	cdeChan := make(chan wl.CallbackDoneEvent)
	cdeHandler := doner{cdeChan}
	callback.AddDoneHandler(cdeHandler)

	sceChan := make(chan wl.SeatCapabilitiesEvent)
	sceHandler := seatcap{sceChan}
	d.seat.AddCapabilitiesHandler(sceHandler)

loop:
	for {
		select {
		case ev := <-sceChan:
			if (ev.Capabilities & wl.SeatCapabilityPointer) != 0 {
				pointer, err := d.seat.GetPointer()
				if err != nil {
					return fmt.Errorf("Unable to get Pointer object: %s", err)
				}
				d.pointer = pointer
			}
			if (ev.Capabilities & wl.SeatCapabilityKeyboard) != 0 {
				keyboard, err := d.seat.GetKeyboard()
				if err != nil {
					return fmt.Errorf("Unable to get Keyboard object: %s", err)
				}
				d.keyboard = keyboard
			}
			if (ev.Capabilities & wl.SeatCapabilityTouch) != 0 {
				touch, err := d.seat.GetTouch()
				if err != nil {
					return fmt.Errorf("Unable to get Touch object: %s", err)
				}
				d.touch = touch
			}
		case d.Dispatch() <- struct{}{}:
		case <-cdeChan:
			break loop
		}
	}

	d.seat.RemoveCapabilitiesHandler(sceHandler)
	callback.RemoveDoneHandler(cdeHandler)

	return nil
}

func (d *Display) registerInterface(registry *wl.Registry, ev wl.RegistryGlobalEvent) error {
	fmt.Printf("we discovered an interface: %q\n", ev.Interface)
	switch ev.Interface {
	case "wl_shm":
		ret := wl.NewShm(d.Context())
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, ret)
		if err != nil {
			return fmt.Errorf("Unable to bind Shm interface: %s", err)
		}
		d.shm = ret
	case "wl_compositor":
		ret := wl.NewCompositor(d.Context())
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, ret)
		if err != nil {
			return fmt.Errorf("Unable to bind Compositor interface: %s", err)
		}
		d.compositor = ret
	case "wl_shell":
		ret := wl.NewShell(d.Context())
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, ret)
		if err != nil {
			return fmt.Errorf("Unable to bind Shell interface: %s", err)
		}
		d.shell = ret
	case "wl_seat":
		ret := wl.NewSeat(d.Context())
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, ret)
		if err != nil {
			return fmt.Errorf("Unable to bind Seat interface: %s", err)
		}
		d.seat = ret
	case "wl_data_device_manager":
		ret := wl.NewDataDeviceManager(d.Context())
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, ret)
		if err != nil {
			return fmt.Errorf("Unable to bind DataDeviceManager interface: %s", err)
		}
		d.dataDeviceManager = ret
	case "wl_subcompositor":
		ret := wl.NewSubcompositor(d.Context())
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, ret)
		if err != nil {
			return fmt.Errorf("Unable to bind Subcompositor interface: %s", err)
		}
		d.subCompositor = ret
	case "zxdg_shell_v6":
		ret := xdg.NewWmBase(d.Context())
		err := registry.Bind(ev.Name, ev.Interface, ev.Version, ret)
		if err != nil {
			return fmt.Errorf("Unable to bind Subcompositor interface: %s", err)
		}
		d.wmBase = ret
		d.wmBase.AddPingHandler(d)
	}
	return nil
}

func (d *Display) HandleDisplayError(ev wl.DisplayErrorEvent) {
	log.Fatalf("Display Error Event: %d - %s - %d", ev.ObjectId.Id(), ev.Message, ev.Code)
}

func (d *Display) newBuffer(width, height, stride int32) (*wl.Buffer, []byte, error) {
	size := stride * height

	file, err := TempFile(int64(size))
	if err != nil {
		return nil, nil, fmt.Errorf("TempFile failed: %s", err)
	}
	defer file.Close()

	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, nil, fmt.Errorf("syscall.Mmap failed: %s", err)
	}

	pool, err := d.shm.CreatePool(file.Fd(), size)
	if err != nil {
		return nil, nil, fmt.Errorf("Shm.CreatePool failed: %s", err)
	}
	defer pool.Destroy()

	buf, err := pool.CreateBuffer(0, width, height, stride, wl.ShmFormatArgb8888)
	if err != nil {
		return nil, nil, fmt.Errorf("Pool.CreateBuffer failed : %s", err)
	}

	return buf, data, nil
}

func (d *Display) registerWindow(w *Window) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.windows = append(d.windows, w)
}

func (d *Display) unregisterWindow(w *Window) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, _w := range d.windows {
		if _w == w {
			d.windows = append(d.windows[:i], d.windows[i+1:]...)
			break
		}
	}
}

// TODO
func (d *Display) FindWindow() *Window {
	return nil
}

func (d *Display) checkGlobalsRegistered() error {
	if d.seat == nil {
		return fmt.Errorf("Seat is not registered")
	}

	if d.compositor == nil {
		return fmt.Errorf("Compositor is not registered")
	}

	if d.shm == nil {
		return fmt.Errorf("Shm is not registered")
	}

	if d.shell == nil {
		return fmt.Errorf("Shell is not registered")
	}

	if d.dataDeviceManager == nil {
		return fmt.Errorf("DataDeviceManager is not registered")
	}

	return nil
}

func (d *Display) Keyboard() *wl.Keyboard {
	return d.keyboard
}

func (d *Display) Pointer() *wl.Pointer {
	return d.pointer
}

func (d *Display) Touch() *wl.Touch {
	return d.touch
}

func (d *Display) checkInputsRegistered() error {
	if d.keyboard == nil {
		return fmt.Errorf("Keyboard is not registered")
	}

	if d.pointer == nil {
		return fmt.Errorf("Pointer is not registered")
	}

	return nil
}

func (d *Display) HandleWmBasePing(ev xdg.WmBasePingEvent) {
	fmt.Printf("ping <%d>\n", ev.Serial)
	d.wmBase.Pong(ev.Serial)
}
