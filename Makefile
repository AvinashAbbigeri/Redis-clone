run: build
	@./bin/Redis

build:
	@go build -o bin/Redis -buildvcs=false