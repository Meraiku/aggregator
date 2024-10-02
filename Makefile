build:
	@go build -o ./.bin/gator ./cmd/gator/

run:build
	@./.bin/gator
