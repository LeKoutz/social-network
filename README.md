# forum

This project can be used to host a forum on the web. It serves via the HTTP
protocol and uses a sqlite3 database to store information.

It supports:

- categorized posts,
- comments on posts,
- likes and dislikes per post/comment,
- user login and registration.

## Build

To build the project you can use the command:
```
export CGO_ENABLED=1; go build -o forum cmd/forum/main.go
```

or

```
make
```
## `docker`
Start it with:
```
docker build -t forum-app .
docker run -p [port]:8080 forum-app
```

## Usage

### `forum`
Start it with:
```
./bin/forum [ip] [port]
```
or
```
./bin/forum [port]
```
or
```
./bin/forum
```
and navigate to `http://[ip]:[port]`

### `mockgen`

`mockgen` can be used to generate some sample data for your DB. Use it by
running:
```
./bin/mockgen
```

### `hashgen`

`hashgen` can be used to generate hashes compatible with the ones stored in the
database. If you want to manually add a user inside the database, you can use:
```
./bin/hashgen <PASSWORD>
```

This will return you hash, which you can enter in your custom query. After
creating the user, you'd be able to login with the password you used to generate
the hash.

## Run tests

To run the tests you can use the command:
```
go test -v ./...
```

or

```
make test
```
