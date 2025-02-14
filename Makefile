GO_ARCH ?= amd64
GO_OS ?= linux
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -gcflags "all=-N -l" -ldflags '-extldflags "-static"' -ldflags "-X github.com/siredmar/elwainthesun/cmd/controller/cmd.Version=$(VERSION)"

build:
	CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=${GO_ARCH} go build $(GO_LDFLAGS) -o bin/controller cmd/controller/main.go 

docker:
	docker build -t siredmar/elwainthesun:$(VERSION) .

push:
	docker push siredmar/elwainthesun:$(VERSION)
