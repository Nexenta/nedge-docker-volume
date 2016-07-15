NDVOL_EXE = ndvol
FLAGS = -v

NEDGE_DEST=/opt/nedge/sbin
NEDGE_ETC = /opt/nedge/etc/ccow

GO_VERSION = 1.6
GO_INSTALL = /usr/lib/go-$(GO_VERSION)
GO = $(GO_INSTALL)/bin/go

all: $(NDVOL_EXE)

GO_FILES = src/github.com/Nexenta/nedge-docker-volume/ndvol/ndvol.go \
	src/github.com/Nexenta/nedge-docker-volume/ndvol/ndvolcli/ndvolcli.go \
	src/github.com/Nexenta/nedge-docker-volume/ndvol/ndvolcli/volumecli.go \
	src/github.com/Nexenta/nedge-docker-volume/ndvol/ndvolcli/daemoncli.go \
	src/github.com/Nexenta/nedge-docker-volume/ndvol/daemon/daemon.go \
	src/github.com/Nexenta/nedge-docker-volume/ndvol/daemon/driver.go \
	src/github.com/Nexenta/nedge-docker-volume/ndvol/ndvolapi/ndvolapi.go

$(GO_FILES): setup

deps: setup
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get github.com/docker/go-plugins-helpers/volume
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get github.com/codegangsta/cli
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get github.com/Sirupsen/logrus
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get github.com/coreos/go-systemd/util
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get github.com/opencontainers/runc/libcontainer/user
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get golang.org/x/net/proxy


$(NDVOL_EXE): $(GO_FILES)
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) install github.com/Nexenta/nedge-docker-volume/ndvol

build:
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) build $(FLAGS) github.com/Nexenta/nedge-docker-volume/ndvol


install: $(NDVOL_EXE)
	cp -f bin/$(NDVOL_EXE) $(NEDGE_DEST)
	cp -f src/github.com/Nexenta/nedge-docker-volume/ndvol/daemon/ndvol.json $(NEDGE_ETC)

setup: 
	mkdir -p src/github.com/Nexenta/nedge-docker-volume/ 
	cp -R ndvol/ src/github.com/Nexenta/nedge-docker-volume/ndvol 

lint:
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get -v github.com/golang/lint/golint
	for file in $$(find . -name '*.go' | grep -v vendor | grep -v '\.pb\.go' | grep -v '\.pb\.gw\.go'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

clean:
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) clean

clobber:
	rm -rf src/github.com/Nexenta/nedge-docker-volume
	rm -rf bin/ pkg/

