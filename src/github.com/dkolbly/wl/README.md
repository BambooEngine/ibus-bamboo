# wl

A wayland protocol implementation in pure Go.

This is a Go implementation of the Wayland protocol.  The protocol
files themselves (`client.go` and `xdg/shell.go`) are built using the
tool in `github.com/dkolbly/wl-scanner` from the XML protocol
specification files.

To test:
```
go get github.com/dkolbly/wl/ui/examples/img  
```

then:
```
$GOPATH/bin/img $GOPATH/src/github.com/dkolbly/wl/ui/examples/img/bsd_daemon.jpg
```

This is a hobby project, forked from a hobby project, `github.com/sternix/wl`.


## Desktops

The image program (`img`) works in both weston and in Ubuntu Gnome in wayland mode.
