VERSION := v0.0.7

all: build

build: go-mod
	@go build -installsuffix bin -ldflags="-w -s" github.com/goweiwen/eir/cmd/eir

build-linux: go-mod
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix bin -ldflags="-w -s" github.com/goweiwen/eir/cmd/eir

go-mod:
	@go mod download

docker: docker-build docker-push

docker-build: build
	@docker build . -t goweiwen/eir:$(VERSION)
	@docker tag goweiwen/eir:$(VERSION) goweiwen/eir:latest

docker-push:
	@docker push goweiwen/eir:$(VERSION)
	@docker push goweiwen/eir:latest

