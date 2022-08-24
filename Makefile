.PHONY: build
build:
	go build .

.PHONY: test
test:
	gotestsum ./...

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: run
run:
	go run . serve


.PHONY: install-tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.48.0
	go install github.com/spf13/cobra-cli@latest
	go install gotest.tools/gotestsum@latest