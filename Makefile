.PHONY: test test-update build-css

test:
	go test ./...

test-update:
	go test ./... -update

build-css:
	cd views/styles && npx tailwindcss -i ./input.css -o ../../public/styles.css --watch
