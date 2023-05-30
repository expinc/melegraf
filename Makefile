export CGO_ENABLED=0

test:
	go test ./... -count=1
