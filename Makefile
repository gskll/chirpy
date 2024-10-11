build:
	@go build -o bin/chirpy ./cmd/chirpy/

run: build
	@./bin/chirpy

test:
	@go test -v ./...
