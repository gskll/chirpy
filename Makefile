build:
	@go build -o bin/chirpy ./cmd/chirpy/

run: build
	@./bin/chirpy

test:
	@go test -v ./...

install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
