
all: build

build:
	go get github.com/rakyll/statik
	statik -src=stacks
	go build -o ctrl cmd/ctrl/main.go
