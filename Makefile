GOLINT := golangci-lint

all: dep lint test bench_test

dep:
	go mod tidy
	go mod download

dep-update:
	go get -t -u ./...

test:
	go test -cover -race -v ./...

bench_test:
	go test -bench=.

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.56.2

lint: dep check-lint ## Lint the files local env
	$(GOLINT) run --timeout=5m -c .golangci.yml