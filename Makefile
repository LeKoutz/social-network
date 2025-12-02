all: forum

test:
	export CGO_ENABLED=1; go test -v ./...

forum:
	export CGO_ENABLED=1; go build -o forum cmd/forum/main.go
