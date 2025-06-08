run : build
	@./bin/redis-go --listenAddr :5432

build:
	@go build -o bin/redis-go