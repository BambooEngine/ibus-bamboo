package wl

import (
	"errors"
	"net"
	"syscall"
)

type Request struct {
	pid    ProxyId
	opcode uint32
	data   []byte
	oob    []byte
}

func (context *Context) SendRequest(proxy Proxy, opcode uint32, args ...interface{}) (err error) {
	req := Request{
		pid:    proxy.Id(),
		opcode: opcode,
	}

	for _, arg := range args {
		req.Write(arg)
	}

	return writeRequest(context.conn, req)
}

func (r *Request) Write(arg interface{}) {
	switch t := arg.(type) {
	case Proxy:
		r.PutProxy(t)
	case uint32:
		r.PutUint32(t)
	case int32:
		r.PutInt32(t)
	case float32:
		r.PutFloat32(t)
	case string:
		r.PutString(t)
	case []int32:
		r.PutArray(t)
	case uintptr:
		r.PutFd(t)
	default:
		panic("Invalid Wayland request parameter type.")
	}
}

func (r *Request) PutUint32(u uint32) {
	buf := bytePool.Take(4)
	order.PutUint32(buf, u)
	r.data = append(r.data, buf...)
}

func (r *Request) PutProxy(p Proxy) {
	r.PutUint32(uint32(p.Id()))
}

func (r *Request) PutInt32(i int32) {
	r.PutUint32(uint32(i))
}

func (r *Request) PutFloat32(f float32) {
	fx := float64ToFixed(float64(f))
	r.PutUint32(uint32(fx))
}

func (r *Request) PutString(s string) {
	tail := 4 - (len(s) & 0x3)
	r.PutUint32(uint32(len(s) + tail))
	r.data = append(r.data, []byte(s)...)
	// if padding required
	if tail > 0 {
		padding := make([]byte, tail)
		r.data = append(r.data, padding...)
	}
}

func (r *Request) PutArray(a []int32) {
	r.PutUint32(uint32(len(a)))
	for _, e := range a {
		r.PutUint32(uint32(e))
	}
}

func (r *Request) PutFd(fd uintptr) {
	rights := syscall.UnixRights(int(fd))
	r.oob = append(r.oob, rights...)
}

func writeRequest(conn *net.UnixConn, r Request) error {
	var header []byte
	// calculate message total size
	size := uint32(len(r.data) + 8)
	buf := make([]byte, 4)
	order.PutUint32(buf, uint32(r.pid))
	header = append(header, buf...)
	order.PutUint32(buf, uint32(size<<16|r.opcode&0x0000ffff))
	header = append(header, buf...)

	d, c, err := conn.WriteMsgUnix(append(header, r.data...), r.oob, nil)
	if err != nil {
		return err
	}
	if c != len(r.oob) || d != (len(header)+len(r.data)) {
		return errors.New("WriteMsgUnix failed")
	}
	bytePool.Give(r.data)

	return nil
}
