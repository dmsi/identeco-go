default: httpserver

all: awslambda

awslambda:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/getjwks ./cmd/awslambda/handlers/jwksets
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/login ./cmd/awslambda/handlers/login
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/register ./cmd/awslambda/handlers/register
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/rotatekeys ./cmd/awslambda/handlers/rotatekeys
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/refresh ./cmd/awslambda/handlers/refresh

httpserver:
	go build -o ./bin/identeco-http ./cmd/httpserver

clean:
	rm bin/*

test:
	go test ./...

run:
	@go run .