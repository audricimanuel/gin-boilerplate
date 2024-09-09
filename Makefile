run_http:
	go run cmd/http/*.go

swag:
	swag init -g cmd/http/main.go ./docs
