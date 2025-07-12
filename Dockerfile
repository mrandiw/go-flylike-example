# Stage 1: Builder
FROM golang:1.22-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go mod tidy
RUN go build -o server .

# Stage 2: Minimal image
FROM alpine:latest

WORKDIR /app

# Only copy the compiled binary
COPY --from=builder /app/server .

EXPOSE 9090

CMD ["./server"]
