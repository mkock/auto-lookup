BINARY := auto-lookup

.PHONY: darwin
darwin:
	mkdir -p release
	GOOS=darwin GOARCH=amd64 go build -o release/$(BINARY)-darwin-amd64 cmd/auto-lookup/main.go

.PHONY: linux
linux:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/$(BINARY)-linux-amd64 cmd/auto-lookup/main.go

.PHONY: clean
clean:
	rm -rf release/*
