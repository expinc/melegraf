export CGO_ENABLED=0

build:
	go build -o gen/melegraf main.go

test:
	go test ./... -count=1
