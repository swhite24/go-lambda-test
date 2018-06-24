

build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler cmd/handler.go

deploy: build
	serverless deploy

deploy-quick: build
	serverless deploy function -f handler
