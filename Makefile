run: install
	@terra-map .

install:
	go install

test:
	go test

.PHONY: run install test
