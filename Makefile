export CGO_ENABLED=0

.PHONY:build
build:
	go build -o gen/melegraf main.go

.PHONY:test
test:
	go test ./... -count=1
