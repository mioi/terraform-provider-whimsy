default: build

.PHONY: build
build:
	go build -v .

.PHONY: install
install: build
	go install -v .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	cd tools; go generate ./...

.PHONY: fmt
fmt:
	gofmt -s -w -e .

.PHONY: test
test:
	go test -v -cover ./internal/provider/

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover ./internal/provider/