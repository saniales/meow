build: cmd/*
	go build -trimpath -o meow

cross-compile: cmd/* pkg/* main.go
	./cross-compile.sh

test:
	go test -v -race -cover ./...