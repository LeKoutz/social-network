# forum

## Build

To build the project you can use the command:
```
export CGO_ENABLED=1; go build -o forum cmd/forum/main.go
```

or

```
make
```

## Usage

Start it with:
```
./forum [ip] [port]
```
or
```
./forum [port]
```
or
```
./forum
```
and navigate to the URL that it returned.

## Run tests

To run the tests you can use the command:
```
go test -v ./...
```

or

```
make test
```
