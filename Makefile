BINARY_NAME := devmux
VERSION := $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: all build clean install

all: build

build:
	go build $(LDFLAGS) -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

install:
	export GOBIN=~/.local/bin/ && go install $(LDFLAGS)
