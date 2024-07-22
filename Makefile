build:
	go build -o ./bin/neko .

run:build
	./bin/neko

tidy:
	go mod tidy

test:
	go test -v
