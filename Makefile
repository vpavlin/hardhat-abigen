OUTPUT=hardhat-abigen
VERSION=0.0.1

build:
	go build -o $(OUTPUT) main.go

release: build
	tar czf hardhat-abigen.$(VERSION).$(shell uname -i).$(shell uname -s | tr '[:upper:]' '[:lower:]').tgz $(OUTPUT)