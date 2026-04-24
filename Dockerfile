FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /gabaa-bot cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /gabaa-bot .

# Copy configuration and migration schemas
COPY --from=builder /app/config/ ./config/
COPY --from=builder /app/internal/constant/query/schemas/ ./internal/constant/query/schemas/

# Expose port (default 8085 as per config)
EXPOSE 8085

CMD ["./gabaa-bot"]
