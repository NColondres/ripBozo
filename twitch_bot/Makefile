.SILENT:

.PHONY: build run clean

build: clean
	go build -o bin/ripBozo

dev: build
	scripts/dev-run.sh

clean:
	rm -rf bin
