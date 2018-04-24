install:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen
	dep ensure

test:
	go vet $(shell go list ./... | grep -v vendor)
	go test -race $(shell go list ./... | grep -v vendor)

mock:
	sh script/mock.sh