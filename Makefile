.PHONY: run seed build build-all clean

run:
	go run ./cmd/streamflix

seed:
	go run ./cmd/seed

build:
	go build -o dist/streamflix ./cmd/streamflix

# Single static binaries, no runtime dependencies (CGO_ENABLED=0), including
# genuine 32-bit targets — this is the whole point of the Go rewrite.
build-all:
	GOOS=windows GOARCH=386   CGO_ENABLED=0 go build -o dist/streamflix-windows-386.exe   ./cmd/streamflix
	GOOS=linux   GOARCH=386   CGO_ENABLED=0 go build -o dist/streamflix-linux-386         ./cmd/streamflix
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/streamflix-windows-amd64.exe ./cmd/streamflix
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o dist/streamflix-linux-amd64       ./cmd/streamflix
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -o dist/streamflix-darwin-arm64      ./cmd/streamflix

clean:
	rm -rf dist
