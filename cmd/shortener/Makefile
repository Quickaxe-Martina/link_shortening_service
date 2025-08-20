start:
	@go fmt .
	@go run main.go

test:
	@go fmt .
	@go test -v -count 1

y_test:
	@./shortenertest -test.v -test.run=^TestIteration1$ -binary-path=./shortener
	@./shortenertest -test.v -test.run=^TestIteration2$ -binary-path=./shortener -source-path=.

build:
	@go fmt .
	@go build -o shortener *.go