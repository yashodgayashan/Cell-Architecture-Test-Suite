# Dockerfile for Go services
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /analytics-frontend main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /analytics-frontend .
EXPOSE 8080
CMD ["./analytics-frontend"]