NDVOL_EXE = ndvol
FLAGS = -v

NEDGE_DEST=/opt/nedge/sbin

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
	GOPATH=$(shell pwd) go get github.com/docker/go-plugins-helpers/volume
	GOPATH=$(shell pwd) go get github.com/codegangsta/cli
	GOPATH=$(shell pwd) go get github.com/Sirupsen/logrus
	GOPATH=$(shell pwd) go get github.com/coreos/go-systemd/util
	GOPATH=$(shell pwd) go get github.com/opencontainers/runc/libcontainer/user
	GOPATH=$(shell pwd) go get golang.org/x/net/proxy


$(NDVOL_EXE): $(GO_FILES)
	GOPATH=$(shell pwd) go install github.com/Nexenta/nedge-docker-volume/ndvol

build:
	GOPATH=$(shell pwd) go build $(FLAGS) github.com/Nexenta/nedge-docker-volume/ndvol


install: $(NDVOL_EXE)
	cp -f bin/$(NDVOL_EXE) $(NEDGE_DEST)

setup: 
	mkdir -p src/github.com/Nexenta/nedge-docker-volume/ 
	cp -R ndvol/ src/github.com/Nexenta/nedge-docker-volume/ndvol 

lint:
	GOPATH=$(shell pwd) go get -v github.com/golang/lint/golint
	for file in $$(find . -name '*.go' | grep -v vendor | grep -v '\.pb\.go' | grep -v '\.pb\.gw\.go'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

clean:
	GOPATH=$(shell pwd) go clean

clobber:
	rm -rf src/github.com/Nexenta/nedge-docker-volume
	rm -rf bin/ pkg/

