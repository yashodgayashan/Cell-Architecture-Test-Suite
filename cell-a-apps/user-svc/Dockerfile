# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
# No dependencies to download yet

COPY . .

# Build the binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-svc .

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /

# Copy the static binary from the builder stage
COPY --from=builder /app/user-svc /user-svc

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/user-svc"]