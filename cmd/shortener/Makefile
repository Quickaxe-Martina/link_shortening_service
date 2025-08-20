start:
	@go fmt .
	@go run main.go

test:
	@go fmt .
	@go test -v -count 1
	@./shortenertest -test.v -test.run=^TestIteration1$ -binary-path=./shortener

build:
	@go fmt .
	@go build -o shortener *.go