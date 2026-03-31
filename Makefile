.PHONY: build clean install lint

BINARY = compartment

build: clean
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY)

lint:
	go vet ./...
	go fmt ./...

install: build
	cp $(BINARY) $(HOME)/.local/bin
