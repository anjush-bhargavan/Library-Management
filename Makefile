run:
		go run main.go

test:
		golint ./...

deps:
		go mod tidy