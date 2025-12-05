all: forum

test:
	export CGO_ENABLED=1; go test -v ./...

forum:
	export CGO_ENABLED=1; go build -o forum cmd/forum/main.go

hashgen:
	export CGO_ENABLED=1; go build -o hashgen cmd/hashgen/main.go

mockgen:
	export CGO_ENABLED=1; go build -o mockgen cmd/mockgen/main.go
