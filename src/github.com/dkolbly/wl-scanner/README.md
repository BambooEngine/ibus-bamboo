wl-scanner
==========

The `wl-scanner` project is designed parse and generate Go client code
from a wayland protocol file.

In its base form, it produces the `client.go` code in
`github.com/dkolbly/wl` from the canonical definition of the wayland
protocol at
https://cgit.freedesktop.org/wayland/wayland/plain/protocol/wayland.xml

It is similar in concept to the wayland-scanner tool which was
developed for generating the C client library.

This is a hobby project intended to help people understand the
wayland.xml protocol file and the generation go code, as well as to
help build client libraries for protocols around Wayland.

## Usage

```
go get github.com/dkolbly/wl-scanner

# generate a client for the base protocol
wl-scanner -source https://cgit.freedesktop.org/wayland/wayland/plain/protocol/wayland.xml \
           -output $GOPATH/src/github.com/dkolbly/wl/client.go

# generate a client for the xdg-shell protocol
wl-scanner -pkg xdg \
           -source https://raw.githubusercontent.com/wayland-project/wayland-protocols/master/stable/xdg-shell/xdg-shell.xml \
           -output $GOPATH/src/github.com/dkolbly/wl/xdg/shell.go
```

### Unstable Protocols

Unstable protocols can be generated in a way that makes it relatively
easy to upgrade to newer versions of the protocol, under the assumption
that changes are relatively small.  This is done by stripping the "_vN" suffix,
as for example:

```
wl-scanner -pkg zxdg \
           -source https://raw.githubusercontent.com/wayland-project/wayland-protocols/master/unstable/xdg-shell/xdg-shell-unstable-v6.xml \
           -output $GOPATH/src/github.com/dkolbly/wl/xdg-unstable-v6/shell.go
```

The `zxdg` package name is preserved, but this can be imported in the client application
as:

```
import (
	xdg "github.com/dkolbly/wl/xdg-unstable-v6"
)
```

and then most code will hopefully be able to port from one major
change to another.  Obviously, mileage will vary depending on the
extent of the change.


