all:
	go build "./..."

test:
	.githooks/pre-commit

run-tokenizer:
	go run cmd/toks-shell/main.go

run-parser:
	go run cmd/pars-shell/main.go

clean:
	go clean
