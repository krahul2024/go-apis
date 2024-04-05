build: 
	@go build -o bin/first

run: build 
	@./bin/first 
