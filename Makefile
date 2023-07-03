# TODO improve!
first:
	@go run .

all:
	# go build -ldflags="-s -w" -v -o ./bin/getjwks ./pkg/handlers/getjwks
	# go build -ldflags="-s -w" -v -o ./bin/rotatekeys ./pkg/handlers/rotatekeys
	# go build -ldflags="-s -w" -v -o ./bin/register ./pkg/handlers/register
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/getjwks ./pkg/handlers/getjwks
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/login ./pkg/handlers/login
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/register ./pkg/handlers/register
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/rotateKeys ./pkg/handlers/rotateKeys
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./bin/refresh ./pkg/handlers/refresh

clean:
	rm bin/*

test:
	go test ./...
