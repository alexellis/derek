all: test build
build:
	go build -o derek

test:
	go test -cover
