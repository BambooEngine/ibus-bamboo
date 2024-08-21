# Ibus bamboo
ibus_e_name="ibus-engine-bamboo"
version="0.8.4"

# Go flags
LDFLAGS="-w -s -X main.Version=${version}"

# Go flags for FreeBSD
CGO_CFLAGS_FREEBSD="-I/usr/local/include -std=gnu99"
CGO_LDFLAGS_FREEBSD="-L/usr/local/lib -lX11 -lXtst -pthread"

# Functions
build_linux() {
	CGO_ENABLED=1 \
	go build -o ${ibus_e_name} -ldflags "${LDFLAGS}" -mod=vendor
}

build_freebsd() {
	CGO_ENABLED=1 \
	CGO_CFLAGS=${CGO_CFLAGS_FREEBSD} \
	CGO_LDFLAGS=${CGO_LDFLAGS_FREEBSD} \
	go build -o ${ibus_e_name} -ldflags "${LDFLAGS}" -mod=vendor
}

# Main script
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
	build_linux
elif [[ "$OSTYPE" == "freebsd"* ]]; then
	build_freebsd
else
	echo "Operating system is not supported currently!"
fi
