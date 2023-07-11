GO := $(HOME)/go/bin/go1.20.5

build:
	$(GO) build -o ./bin/neko .

run:build
	./bin/neko
