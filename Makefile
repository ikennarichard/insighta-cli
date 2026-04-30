BINARY_NAME=insighta
MODULE=github.com/ikennarichard/insighta-cli

.PHONY: build install clean

build:
	go build -o $(BINARY_NAME) .

# Installs globally to $GOPATH/bin
install:
	go install .

clean:
	go clean
	rm -f $(BINARY_NAME)

# Cross-compile for all platforms
release:
	GOOS=linux   GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin  GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe .