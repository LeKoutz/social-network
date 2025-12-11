FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1; go build -o bin/forum cmd/forum/main.go
EXPOSE 8080
CMD ["./forum"]