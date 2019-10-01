
all: build

build:
	go get github.com/rakyll/statik
	statik -src=stacktmpl
	go build -o ctrl cmd/ctrlcli/main.go
