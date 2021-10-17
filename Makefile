all: test

clean:
	rm -f pcr

install: prepare
	go install

build: prepare
	go build

test: prepare build
	echo "no tests"

.PHONY: install prepare build test