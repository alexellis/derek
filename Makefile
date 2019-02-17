all: test build
build:
	go build -o derek

docker:
	faas-cli build -f derek.yml

test:
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /build/ | grep -v /template/) -cover
