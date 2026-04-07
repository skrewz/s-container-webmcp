.PHONY: build test clean container-build container-run

build:
	go build -o server ./cmd/server

test:
	go test ./...

clean:
	rm -f server

container-build:
	podman build -t s-container-webmcp .

container-run:
	podman run -e DEBUG -p 3952:3952 s-container-webmcp
