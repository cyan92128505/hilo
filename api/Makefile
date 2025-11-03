test:
	@go test $(shell go list ./... | grep -v /tmp) -cover

wire:
	@go mod download && go mod tidy && wire ./cmd/restful