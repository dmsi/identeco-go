default: httpserver awslambda package

all: awslambda package

# awslambda:
# 	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/jwksets ./cmd/awslambda/jwksets
# 	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/register ./cmd/awslambda/register
# 	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/login ./cmd/awslambda/login
# 	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/refresh ./cmd/awslambda/refresh
# 	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/rotatekeys ./cmd/awslambda/rotatekeys

awslambda:
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/jwksets/bootstrap ./cmd/awslambda/jwksets
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/register/bootstrap ./cmd/awslambda/register
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/login/bootstrap ./cmd/awslambda/login
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/refresh/bootstrap ./cmd/awslambda/refresh
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/rotatekeys/bootstrap ./cmd/awslambda/rotatekeys

package:
	zip -j ./bin/jwksets.zip ./bin/jwksets/bootstrap
	zip -j ./bin/register.zip ./bin/register/bootstrap
	zip -j ./bin/login.zip ./bin/login/bootstrap
	zip -j ./bin/refresh.zip ./bin/refresh/bootstrap
	zip -j ./bin/rotatekeys.zip ./bin/rotatekeys/bootstrap

httpserver:
	go build -o ./bin/identeco-http ./cmd/httpserver

clean:
	rm bin/*

test:
	go test ./...
