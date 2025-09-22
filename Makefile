BINARY_NAME=goTickersNotifier
VERSION=1.0.0

build:
	go build -o bin/$(BINARY_NAME) .

build-all:
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64-v$(VERSION) .
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64-v$(VERSION) .
    
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64-v$(VERSION).exe .
	GOOS=windows GOARCH=arm64 go build -o bin/$(BINARY_NAME)-windows-arm64-v$(VERSION).exe .
    
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64-v$(VERSION) .
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64-v$(VERSION) .

clean:
	rm -rf bin/

.PHONY: build build-all clean
