# Go flags for FreeBSD
CGO_CFLAGS_FREEBSD="-I/usr/local/include -std=gnu99"
CGO_LDFLAGS_FREEBSD="-L/usr/local/lib -lX11 -lXtst -pthread"

# Functions
test_linux() {
        CGO_ENABLED=1 \
        go test . -mod=vendor
}

test_freebsd() {
        CGO_ENABLED=1 \
        CGO_CFLAGS=${CGO_CFLAGS_FREEBSD} \
        CGO_LDFLAGS=${CGO_LDFLAGS_FREEBSD} \
        go test . -mod=vendor
}

# Main script
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        test_linux
elif [[ "$OSTYPE" == "freebsd"* ]]; then
        test_freebsd
else
        echo "Operating system is not supported currently!"
fi
