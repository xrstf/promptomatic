export CGO_ENABLED ?= 0
export GOFLAGS ?= -mod=readonly -trimpath
export GO111MODULE = on
GOTOOLFLAGS ?= -v -ldflags '-w -s -extldflags '-static''

.PHONY: default
default: build

.PHONY: build
build:
	go build $(GOTOOLFLAGS)

.PHONY: install
install:
	go install

.PHONY: sysinstall
sysinstall: build
	sudo mv promptomatic /usr/local/bin/
