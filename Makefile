.PHONY: dev build generate run

dev:
	air

generate:
	~/go/bin/templ generate

build: generate
	go build -o tmp/main ./cmd/server

run: generate
	go run ./cmd/server
