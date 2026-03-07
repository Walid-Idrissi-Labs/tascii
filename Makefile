BINARY := tascii


.PHONY: build
build:
	go build -o $(BINARY) .


.PHONY: install
install:
	go install .


.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	rm -f $(BINARY)


.PHONY: deps
deps:
	go mod tidy


.PHONY: release-dry
release-dry:
	goreleaser release --snapshot --clean
