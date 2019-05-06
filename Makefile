NEDGE_DEST = $(DESTDIR)/opt/nedge/sbin
NEDGE_ETC = $(DESTDIR)/opt/nedge/etc/ccow
NDVOL_EXE = ndvol

ifeq ($(GOPATH),)
GOPATH = $(shell pwd)
endif

build: 
	# docker/go-plugins-helpers
	GOPATH=$(GOPATH) go get -d -v github.com/docker/go-plugins-helpers/volume
	cd $(GOPATH)/src/github.com/docker/go-plugins-helpers; git checkout d7fc7d0
	# opencontainers/runc
	GOPATH=$(GOPATH) go get -d -v github.com/opencontainers/runc
	cd $(GOPATH)/src/github.com/opencontainers/runc; git checkout aada2af
	# docker/go-connections
	GOPATH=$(GOPATH) go get -d -v github.com/docker/go-connections
	cd $(GOPATH)/src/github.com/docker/go-connections; git checkout acbe915
	GOPATH=$(GOPATH) go get -v github.com/Nexenta/nedge-docker-volume/...

lint:
	GOPATH=$(GOPATH) GOROOT=$(GO_INSTALL) $(GO) get -v github.com/golang/lint/golint
	for file in $$(find . -name '*.go' | grep -v vendor | grep -v '\.pb\.go' | grep -v '\.pb\.gw\.go'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

install:
	cp -n ndvol/daemon/ndvol.json $(NEDGE_ETC)
	cp -f bin/$(NDVOL_EXE) $(NEDGE_DEST)

uninstall:
	rm -f $(NEDGE_ETC)/ndvol.json
	rm -f $(NEDGE_DEST)/ndvol

clean:
	GOPATH=$(GOPATH) go clean github.com/Nexenta/nedge-docker-volume
