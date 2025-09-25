start:
	@go run ./cmd/shortener/main.go -d="host=127.0.0.1 port=5469 user=test password=test dbname=test sslmode=disable"

build:
	@go build -o ./cmd/shortener/shortener ./cmd/shortener/main.go