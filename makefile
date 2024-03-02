.PHONY: build

build:
	go build -o ./bin/spotle-api ./cmd/.

.PHONY: run

run: build
	./bin/spotle-api