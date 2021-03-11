package wl

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func init() {
	log.SetFlags(0)
}

type Context struct {
	mu           sync.RWMutex
	conn         *net.UnixConn
	currentId    ProxyId
	objects      map[ProxyId]Proxy
	dispatchChan chan struct{}
	exitChan     chan struct{}
}

func (ctx *Context) Register(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.currentId += 1
	proxy.SetId(ctx.currentId)
	proxy.SetContext(ctx)
	ctx.objects[ctx.currentId] = proxy
}

func (ctx *Context) lookupProxy(id ProxyId) Proxy {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	proxy, ok := ctx.objects[id]
	if !ok {
		return nil
	}
	return proxy
}

func (ctx *Context) unregister(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	delete(ctx.objects, proxy.Id())
}

func (c *Context) Close() {
	c.conn.Close()
	c.exitChan <- struct{}{}
	close(c.dispatchChan)

}

func (c *Context) Dispatch() chan<- struct{} {
	return c.dispatchChan
}

func Connect(addr string) (ret *Display, err error) {
	runtime_dir := os.Getenv("XDG_RUNTIME_DIR")
	if runtime_dir == "" {
		return nil, errors.New("XDG_RUNTIME_DIR not set in the environment.")
	}
	if addr == "" {
		addr = os.Getenv("WAYLAND_DISPLAY")
	}
	if addr == "" {
		addr = "wayland-0"
	}
	addr = runtime_dir + "/" + addr
	c := new(Context)
	c.objects = make(map[ProxyId]Proxy)
	c.currentId = 0
	c.dispatchChan = make(chan struct{})
	c.exitChan = make(chan struct{})
	c.conn, err = net.DialUnix("unix", nil, &net.UnixAddr{Name: addr, Net: "unix"})
	if err != nil {
		return nil, err
	}
	c.conn.SetReadDeadline(time.Time{})
	//dispatch events in separate gorutine
	go c.run()
	return NewDisplay(c), nil
}

func (c *Context) run() {
	ctx := context.Background()

loop:
	for {
		select {
		case <-c.dispatchChan:
			ev, err := c.readEvent()
			if err != nil {
				if err == io.EOF {
					// connection closed
					break loop

				}

				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					log.Print("Timeout Error")
					continue
				}

				log.Fatal(err)
			}

			proxy := c.lookupProxy(ev.pid)
			if proxy != nil {
				if dispatcher, ok := proxy.(Dispatcher); ok {
					dispatcher.Dispatch(ctx, ev)
					bytePool.Give(ev.data)
				} else {
					log.Print("Not dispatched")
				}
			} else {
				log.Print("Proxy NULL")
			}

		case <-c.exitChan:
			break loop
		}
	}
}
