FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o bin/forum cmd/forum/main.go
RUN CGO_ENABLED=1 go build -o bin/mockgen cmd/mockgen/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/forum ./app/forum
COPY --from=builder /app/bin/mockgen ./app/mockgen
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations
RUN mkdir -p /app/data
EXPOSE 8080
RUN ./app/mockgen /app/data/db.db
CMD ["./app/forum", "--db-path", "/app/data/db.db", "0.0.0.0", "8080"]
