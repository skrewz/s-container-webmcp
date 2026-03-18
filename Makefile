.PHONY: build test clean container-build container-run

build:
	go build -o server ./cmd/server

test:
	go test ./...

clean:
	rm -f server

container-build:
	podman build -t web-search-mcp .

container-run:
	podman run -p 3952:3952 web-search-mcp
