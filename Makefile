.PHONY: all build clean test lint vet check cross-compile

# Default target
all: build

# Build the main binary
build:
	go build -o fp2lm ./cmd/fp2lm

# Clean built files
clean:
	rm -f fp2lm *.bin
	rm -f TestOutput.csv TestFixed.csv TestFixedASL.csv FinalTest.csv
	rm -rf dist/

# Run tests
test:
	go test -v ./...

# Run linters
lint: vet
	@echo "Linting code..."
	go install honnef.co/go/tools/cmd/staticcheck@latest
	$(shell go env GOPATH)/bin/staticcheck ./...
	go install golang.org/x/vuln/cmd/govulncheck@latest
	$(shell go env GOPATH)/bin/govulncheck ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run all checks
check: test lint

# Run tests with race detection
test-race:
	go test -race -v ./...

# Cross-compile binaries for different platforms
cross-compile:
	@echo "Cross-compiling for multiple platforms..."
	mkdir -p dist
	# Linux (amd64)
	GOOS=linux GOARCH=amd64 go build -o dist/fp2lm-linux-amd64 ./cmd/fp2lm
	# macOS (amd64)
	GOOS=darwin GOARCH=amd64 go build -o dist/fp2lm-darwin-amd64 ./cmd/fp2lm
	# macOS (arm64)
	GOOS=darwin GOARCH=arm64 go build -o dist/fp2lm-darwin-arm64 ./cmd/fp2lm
	# Windows (amd64)
	GOOS=windows GOARCH=amd64 go build -o dist/fp2lm-windows-amd64.exe ./cmd/fp2lm

# Create release archives
release: cross-compile
	@echo "Creating release archives..."
	cd dist && tar -czf fp2lm-linux-amd64.tar.gz fp2lm-linux-amd64
	cd dist && tar -czf fp2lm-darwin-amd64.tar.gz fp2lm-darwin-amd64
	cd dist && tar -czf fp2lm-darwin-arm64.tar.gz fp2lm-darwin-arm64
	cd dist && zip fp2lm-windows-amd64.zip fp2lm-windows-amd64.exe 