FROM golang:1.17.1 AS builder
# Set work directory
WORKDIR /src

# Copy go module files
COPY go.mod go.sum ./

# Download go module dependencies
RUN go mod download all

# Copy source files and build executable
COPY . ./
RUN go build -v -o /bin/app ./*.go

# Production container
FROM debian:bullseye-slim

# Set working directory
WORKDIR /app

# Update and install apt dependencies
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -q - y

# Copy config.json file
COPY config.json ./

# Copy executable from builder
COPY --from=builder /bin/app ./

# Start app
CMD ["./app"]