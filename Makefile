build:
	go build -C ./cmd/books -o ../../bin/books

run:build
	-docker compose up -d
	-GIN_MODE=release ./bin/books

dev:build
	-docker compose up -d
	./bin/books
