build:
	GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o bin/controller cmd/controller/main.go 

docker:
	docker build -t siredmar/elwainthesun:latest .

push:
	docker push siredmar/elwainthesun:latest
