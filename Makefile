GOBIN ?= $$(go env GOPATH)/bin

.PHONY: test test-cover lint cover-html

build:
	go build -o ./limepipes-plugin-bww github.com/tomvodi/limepipes-plugin-bww/cmd/limepipes-plugin-bww

mocks:
	mockery

test:
	go test ./...

test-cover:
	go test ./... -coverprofile cover.out

lint:
	golangci-lint run

cover-html: test-cover
	go tool cover -html=cover.out

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yaml