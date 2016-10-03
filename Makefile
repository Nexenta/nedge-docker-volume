setup: 
	sudo cp ndvol/daemon/ndvol.json /opt/nedge/etc/ccow/
	go get -v github.com/Nexenta/nedge-docker-volume/...

lint:
	go get -v github.com/golang/lint/golint
	for file in $$(find $GOPATH/src/github.com/Nexenta/nedge-docker-volume -name '*.go' | grep -v vendor | grep -v '\.pb\.go' | grep -v '\.pb\.gw\.go'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

clean:
	go clean github.com/Nexenta/nedge-docker-volume
