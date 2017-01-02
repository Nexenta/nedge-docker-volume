NEDGE_DEST = $(DESTDIR)/opt/nedge/sbin
NEDGE_ETC = $(DESTDIR)/opt/nedge/etc/ccow
NDVOL_EXE = ndvol

build: 
	GOPATH=$(shell pwd) go get -v github.com/docker/go-plugins-helpers/volume
	cd src/github.com/docker/go-plugins-helpers/volume; git checkout d7fc7d0
	cd src/github.com/docker/go-connections; git checkout acbe915
	GOPATH=$(shell pwd) go get -v github.com/Nexenta/nedge-docker-volume/...

lint:
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get -v github.com/golang/lint/golint
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
	go clean github.com/Nexenta/nedge-docker-volume
