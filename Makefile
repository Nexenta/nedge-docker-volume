NED_EXE = nedv
FLAGS = -v

all: $(NED_EXE)

GO_FILES = src/github.com/Nexenta/nedge-docker-volume/nedv/nedv.go \
	src/github.com/Nexenta/nedge-docker-volume/nedv/nedcli/nedcli.go \
	src/github.com/Nexenta/nedge-docker-volume/nedv/nedcli/Foo.go \
	src/github.com/Nexenta/nedge-docker-volume/nedv/nedcli/Bar.go \
	src/github.com/Nexenta/nedge-docker-volume/nedv/daemon/daemon.go \
	src/github.com/Nexenta/nedge-docker-volume/nedv/daemon/driver.go \
	src/github.com/Nexenta/nedge-docker-volume/nedv/nedapi/nedapi.go

$(GO_FILES): setup

deps: setup
	GOPATH=$(shell pwd) go get github.com/docker/go-plugins-helpers/volume
	GOPATH=$(shell pwd) go get github.com/codegangsta/cli
	GOPATH=$(shell pwd) go get github.com/Sirupsen/logrus
	GOPATH=$(shell pwd) go get github.com/coreos/go-systemd/util
	GOPATH=$(shell pwd) go get github.com/opencontainers/runc/libcontainer/user
	GOPATH=$(shell pwd) go get golang.org/x/net/proxy


$(NED_EXE): $(GO_FILES)
	GOPATH=$(shell pwd) go build $(FLAGS) github.com/Nexenta/nedge-docker-volume/nedv

setup: 
	mkdir -p src/github.com/Nexenta/nedge-docker-volume/ 
	cp -R ned/ src/github.com/Nexenta/nedge-docker-volume/nedv 

clean:
	rm $(NED_EXE) \
	rm -rf bin/ pkg/


clean_src:
	rm -rf src/github.com/Nexenta/nedge-docker-volume

install:
	go install github.com/Nexenta/nedge-docker-volume/ned
