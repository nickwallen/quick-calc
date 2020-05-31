TARGET = bin
BINARY = qcalc

all: clean build test 

clean:
	go clean "./..."
	rm -rf $(TARGET)

build:
	go fmt ./...
	go vet ./...
	golint "-set_exit_status" ./...
	go build "./..."

test:
	go test "./..."

debug:
	go run cmd/cli/main.go debug

run: build
	go run cmd/cli/main.go

binary:
	mkdir $(TARGET)
	go build -o $(TARGET)/$(BINARY) cmd/qcalc-cli/main.go

samples:
	go run cmd/gen/main.go

