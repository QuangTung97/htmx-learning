.PHONY: test build-css

test:
	go test ./...

build-css:
	cd views/styles && npx tailwindcss -i ./input.css -o ../../public/styles.css
