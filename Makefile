# Makefile — convenience commands for development.
# Run any target with: make <target>  (e.g. make build)

# The binary name.
BINARY := tascii

# Default target: build for your current machine.
.PHONY: build
build:
	go build -o $(BINARY) .

# Install into your Go bin directory so you can run `tascii` directly.
# Your Go bin (~/.local/bin or ~/go/bin) should be in your PATH.
.PHONY: install
install:
	go install .

# Run all tests.
.PHONY: test
test:
	go test ./...

# Remove the built binary.
.PHONY: clean
clean:
	rm -f $(BINARY)

# Download all dependencies listed in go.mod.
.PHONY: deps
deps:
	go mod tidy

# Dry-run GoReleaser to check your .goreleaser.yaml config without publishing.
.PHONY: release-dry
release-dry:
	goreleaser release --snapshot --clean
