TARGET = bin
BINARY = qcalc

all: clean check build test 

clean:
	go clean "./..."
	rm -rf $(TARGET)

check:
	go fmt ./...
	go vet ./...
	golint "-set_exit_status" ./...

build:
	go build "./..."

test:
	go test "./..."

debug:
	go run cmd/cli/main.go debug

run:
	go run cmd/cli/main.go

binary:
	mkdir $(TARGET)
	go build -o $(TARGET)/$(BINARY) cmd/qcalc-cli/main.go

samples:
	go run cmd/gen/main.go

