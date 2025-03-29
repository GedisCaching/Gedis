.DEFAULT_GOAL := run

.PHONY:fmt vet run
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

run: vet
	go run main.go