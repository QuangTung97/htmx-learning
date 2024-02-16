.PHONY: test test-update build-css run install-tools migrate-up

test:
	go test ./...

test-update:
	GOLDEN_UPDATE=1 go test ./...

build-css:
	cd views/styles && npx tailwindcss -i ./input.css -o ../../public/styles.css --watch

run:
	go run cmd/main.go

install-tools:
	go install github.com/matryer/moq
	go install github.com/mgechev/revive

migrate-up:
	go run cmd/migrate/main.go up