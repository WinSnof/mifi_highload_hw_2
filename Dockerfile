FROM golang:1.25.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go-microservice main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go-microservice .
EXPOSE 8080
CMD ["./go-microservice"]