TARGET = bin
BINARY = qcalc

all:
	go build "./..."

test:
	.githooks/pre-commit

debug:
	go run cmd/cli/main.go debug

run:
	go run cmd/cli/main.go

binary:
	mkdir $(TARGET)
	go build -o $(TARGET)/$(BINARY) cmd/qcalc-cli/main.go

samples:
	go run cmd/gen/main.go

clean:
	go clean "./..."
	rm -rf $(TARGET)
