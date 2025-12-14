FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o bin/forum cmd/forum/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/forum ./app/forum
COPY --from=builder /app/templates ./templates
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["./app/forum", "--db-path", "/app/data/db.db", "0.0.0.0", "8080"]