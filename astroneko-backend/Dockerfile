# Build stage
FROM golang:1.24.3-alpine AS builder

# Install git
RUN apk add --no-cache git ca-certificates

# Set build arguments
ARG GITHUB_USERNAME
ARG GITHUB_TOKEN

# Set environment variables for private module
RUN mkdir /user && \
    echo 'root:x:65534:65534:root:/:' > /user/passwd && \
    echo 'root:x:65534:' > /user/group && \
    printf "machine github.com login $GITHUB_USERNAME password $GITHUB_TOKEN" >> /root/.netrc && \
    chmod 600 /root/.netrc

# Set working directory
WORKDIR /app


# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/api

# Final stage
FROM alpine:3.20

WORKDIR /app


# Create configs directory with proper permissions
RUN apk --no-cache add ca-certificates && \
    mkdir -p /app/configs && \
    chmod 755 /app/configs

# Copy the binary from builder
COPY --from=builder /app/app .
COPY --from=builder /app/configs/prod-config.yml /app/configs/

# Verify config file exists and set permissions
RUN ls -la /app/configs/prod-config.yml && chmod 644 /app/configs/prod-config.yml

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"] 