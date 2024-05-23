build: cmd/*
	go build -o meow

cross-compile: cmd/*
	./cross-compile.sh
