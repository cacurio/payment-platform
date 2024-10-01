
build:
	GOOS=linux GOARCH=arm64 go build -o ./cmd/refund/bootstrap  ./cmd/refund/handler.go
	GOOS=linux GOARCH=arm64 go build -o ./cmd/charge/bootstrap  ./cmd/charge/handler.go
	GOOS=linux GOARCH=arm64 go build -o ./cmd/token/bootstrap  ./cmd/token/handler.go
	GOOS=linux GOARCH=arm64 go build -o ./cmd/get-charge/bootstrap  ./cmd/get-charge/handler.go

