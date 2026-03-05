.PHONY: run build clean test deps

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

clean:
	rm -rf bin/

test:
	go test -v ./...

deps:
	go mod download
	go mod tidy

install:
	go install ./...

docker-build:
	docker build -t gin-htmx-template .

docker-run:
	docker run -p 5007:5007 --env-file .env gin-htmx-template
