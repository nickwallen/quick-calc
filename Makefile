all:
	go build "./..."

test:
	.githooks/pre-commit

run-tokenizer:
	go run cmd/calc-cli/main.go tokens

run:
	go run cmd/calc-cli/main.go

clean:
	go clean
