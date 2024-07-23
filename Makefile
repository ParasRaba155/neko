build:
	go build -o ./bin/neko .

run:build
	./bin/neko

tidy:
	go mod tidy

test-verbose:
	go test -v

test:
	go test ./...

format:
	go fmt ./...
