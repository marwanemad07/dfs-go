# this docker file is used for data nodes
# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o datanode_app ./datanode/data_keeper.go

# Stage 2: Run (minimal image)
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/datanode_app .

ENTRYPOINT ["./datanode_app", "docker"]
CMD ["5050"]

