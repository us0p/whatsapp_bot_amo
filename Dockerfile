# ---- Builder stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (needed for some dependencies)
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o app

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port (optional but nice for clarity)
EXPOSE 8000

# Run the app
CMD ["./app"]
