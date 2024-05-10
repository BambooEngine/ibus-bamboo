package wl

import (
	"bytes"
	"fmt"
	"syscall"
)

type Event struct {
	pid    ProxyId
	Opcode uint32
	data   []byte
	scms   []syscall.SocketControlMessage
	off    int
}

func (c *Context) readEvent() (*Event, error) {
	buf := bytePool.Take(8)
	control := bytePool.Take(24)

	n, oobn, _, _, err := c.conn.ReadMsgUnix(buf[:], control)
	if err != nil {
		return nil, err
	}
	if n != 8 {
		return nil, fmt.Errorf("Unable to read message header.")
	}
	ev := new(Event)
	if oobn > 0 {
		if oobn > len(control) {
			return nil, fmt.Errorf("Unsufficient control msg buffer")
		}
		scms, err := syscall.ParseSocketControlMessage(control)
		if err != nil {
			return nil, fmt.Errorf("Control message parse error: %s", err)
		}
		ev.scms = scms
	}

	ev.pid = ProxyId(order.Uint32(buf[0:4]))
	ev.Opcode = uint32(order.Uint16(buf[4:6]))
	size := uint32(order.Uint16(buf[6:8]))

	// subtract 8 bytes from header
	data := bytePool.Take(int(size) - 8)
	n, err = c.conn.Read(data)
	if err != nil {
		return nil, err
	}
	if n != int(size)-8 {
		return nil, fmt.Errorf("Invalid message size.")
	}
	ev.data = data

	bytePool.Give(buf)
	bytePool.Give(control)

	return ev, nil
}

func (ev *Event) FD() uintptr {
	if ev.scms == nil {
		return 0
	}
	fds, err := syscall.ParseUnixRights(&ev.scms[0])
	if err != nil {
		panic("Unable to parse unix rights")
	}
	//TODO is this required
	ev.scms = append(ev.scms, ev.scms[1:]...)
	return uintptr(fds[0])
}

func (ev *Event) Uint32() uint32 {
	buf := ev.next(4)
	if len(buf) != 4 {
		panic("Unable to read unsigned int")
	}
	return order.Uint32(buf)
}

func (ev *Event) Proxy(c *Context) Proxy {
	return c.lookupProxy(ProxyId(ev.Uint32()))
}

func (ev *Event) String() string {
	l := int(ev.Uint32())
	buf := ev.next(l)
	if len(buf) != l {
		panic("Unable to read string")
	}
	ret := string(bytes.TrimRight(buf, "\x00"))
	//padding to 32 bit boundary
	if (l & 0x3) != 0 {
		ev.next(4 - (l & 0x3))
	}
	return ret
}

func (ev *Event) Int32() int32 {
	return int32(ev.Uint32())
}

func (ev *Event) Float32() float32 {
	return float32(fixedToFloat64(ev.Int32()))
}

func (ev *Event) Array() []int32 {
	l := int(ev.Uint32())
	arr := make([]int32, l/4)
	for i := range arr {
		arr[i] = ev.Int32()
	}
	return arr
}

func (ev *Event) next(n int) []byte {
	ret := ev.data[ev.off : ev.off+n]
	ev.off += n
	return ret
}
