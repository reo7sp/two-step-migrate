.PHONY: prepare build clean

all: build

prepare:
	dep ensure
	go get github.com/mitchellh/gox

build:
	gox -osarch="linux/amd64 darwin/amd64"

clean:
	rm -f two-step-migrate*
