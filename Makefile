build: cmd/*
	go build -o meow

cross-compile: cmd/* pkg/* main.go
	./cross-compile.sh
