.PHONY: run
run: hello
	./hello

hello: hello.s
	gcc -o hello hello.s

hello.s: assembra.go
	make imports
	make format
	go run . 3 > hello.s

.PHONY: format
format:
	gofmt -w .

.PHONY: imports
imports:
	goimports -w .
