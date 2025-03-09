run: build
	@./bin/Redis -listenaddr :5001

build:
	@go build -o bin/Redis .
