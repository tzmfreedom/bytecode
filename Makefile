.PHONY: run
run: format
	@go run .

.PHONY: test
test:
	go test ./...

.PHONY: benchmark
benchmark:
	go test -bench=.

.PHONY: format
format:
	@goimports -w .
	@gofmt -w .

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: deps
deps: dep
	dep ensure
