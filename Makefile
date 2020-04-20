all:
	go build "./..."

test:
	.githooks/pre-commit

run-tokenizer:
	go run cmd/toks-shell/main.go

run:
	go run cmd/calc-cli/main.go

clean:
	go clean
