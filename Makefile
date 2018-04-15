test:
	go vet $(shell go list ./... | grep -v vendor)
	go test -race $(shell go list ./... | grep -v vendor)