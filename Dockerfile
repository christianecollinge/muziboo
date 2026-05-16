FROM golang:1.26-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy the go mod files from the api directory
COPY api/go.mod api/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the api source code
COPY api/ ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api

# Final minimal stage
FROM alpine:latest

# Install CA certificates so the Go server can make secure HTTPS requests (like to Supabase)
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /server /server

# The Go code uses runtime.Caller(0) which bakes the build path into the binary.
# We must copy the templates into the exact same path structure they were built from.
COPY --from=builder /app/foundation/web/templates /app/foundation/web/templates

# Expose the port the Go server runs on
EXPOSE 8080

# Start the server
CMD ["/server"]
