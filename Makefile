install:
	go mod download
	go install github.com/securego/gosec/v2/cmd/gosec@latest

build:
	go build -o ui-dev-proxy main.go

test:
	go vet ./...
	gosec ./...
	go test -race -short ./...
	bash -c 'diff -u <(echo -n) <(gofmt -s -d .)'

fmt:
	go fmt ./...
