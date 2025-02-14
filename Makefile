GO_ARCH ?= amd64
GO_OS ?= linux
VERSION ?= $(shell git describe --tags --always --dirty --abbrev=10)
GO_LDFLAGS = -gcflags "all=-N -l" -ldflags '-extldflags "-static"' -ldflags "-X github.com/siredmar/ElwaInTheSun/cmd/controller/cmd.Version=$(VERSION)"

build:
	CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=${GO_ARCH} go build $(GO_LDFLAGS) -o bin/controller cmd/controller/main.go 

nextversion:
	echo $(VERSION)

docker:
	docker build -t siredmar/elwainthesun:$(VERSION) .

push:
	docker push siredmar/elwainthesun:$(VERSION)
