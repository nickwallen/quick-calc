all:
	go build "./..."

test:
	.githooks/pre-commit

debug:
	go run cmd/calc-cli/main.go debug

run:
	go run cmd/calc-cli/main.go

clean:
	go clean
