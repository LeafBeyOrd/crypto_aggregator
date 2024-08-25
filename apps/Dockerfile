# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o /crypto-processor

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the built Go app from the build stage
COPY --from=build /crypto-processor .

# Command to run the binary
CMD ["./crypto-processor"]
