.PHONY: test test-update build-css run

test:
	go test ./...

test-update:
	go test ./... -update

build-css:
	cd views/styles && npx tailwindcss -i ./input.css -o ../../public/styles.css --watch

run:
	go run cmd/main.go