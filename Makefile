build:
	@go build -o bin/easyPay cmd/main.go

run:build
	@./bin/easyPay