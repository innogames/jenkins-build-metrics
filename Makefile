.PHONY: all clean

all: server

server:
	go build -o bin/jenkins-build-metrics main.go

server\:linux: init
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -o bin/jenkins-build-metrics main.go
