default: standalone awslambda package

all: awslambda package

awslambda:
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/jwksets/bootstrap ./cmd/awslambda/jwksets
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/register/bootstrap ./cmd/awslambda/register
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/login/bootstrap ./cmd/awslambda/login
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/refresh/bootstrap ./cmd/awslambda/refresh
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/rotatekeys/bootstrap ./cmd/awslambda/rotatekeys
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags="-s -w" -v -o ./bin/ping/bootstrap ./cmd/awslambda/ping

package:
	zip -j ./bin/jwksets.zip ./bin/jwksets/bootstrap
	zip -j ./bin/register.zip ./bin/register/bootstrap
	zip -j ./bin/login.zip ./bin/login/bootstrap
	zip -j ./bin/refresh.zip ./bin/refresh/bootstrap
	zip -j ./bin/rotatekeys.zip ./bin/rotatekeys/bootstrap
	zip -j ./bin/ping.zip ./bin/ping/bootstrap

standalone:
	go build -o ./bin/identeco-standalone ./cmd/standalone/server.go

docker:
	@echo "Yes, please!"

clean:
	rm -rf bin/*

test:
	go test ./...
