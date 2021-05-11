
all: build

run:
	go get github.com/rakyll/statik
	statik -src=stacktmpl
	go run ./cmd/ctrl/main.go ./cmd/ctrl/logging.go

build-cli:
	go build -o ctrli ./cmd/ctrlcli
