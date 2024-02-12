.PHONY: test test-update build-css run install-tools

test:
	go test ./...

test-update:
	go test ./... -update

build-css:
	cd views/styles && npx tailwindcss -i ./input.css -o ../../public/styles.css --watch

run:
	go run cmd/main.go

install-tools:
	go install github.com/matryer/moq
	go install github.com/mgechev/revive
