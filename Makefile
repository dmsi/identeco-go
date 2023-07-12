default: httpserver awslambda

all: awslambda

awslambda:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/jwksets ./cmd/awslambda/jwksets
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/register ./cmd/awslambda/register
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/login ./cmd/awslambda/login
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/refresh ./cmd/awslambda/refresh
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/rotatekeys ./cmd/awslambda/rotatekeys

httpserver:
	go build -o ./bin/identeco-http ./cmd/httpserver

clean:
	rm bin/*

test:
	go test ./...
