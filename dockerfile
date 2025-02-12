# Use the official Alpine Linux image
FROM golang:alpine
# Install Go and required libraries
RUN apk add --no-cache go sqlite libc6-compat

# Set working directory inside the container
WORKDIR /app
# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go application from the cmd directory
RUN go build -o main ./cmd/main.go

# Copy the database and schema files
COPY internal/database/forum.db ./internal/database/forum.db
COPY internal/database/schema.sql ./internal/database/schema.sql

# Ensure the binary is executable
RUN chmod +x /app/main

# Set environment variables inside the container
ENV PORT=8081
ENV DB_PATH=/app/internal/database/forum.db
ENV SCHEMA_PATH=/app/internal/database/schema.sql

# Expose the application port
EXPOSE 8081

# Run the application
CMD ["./main"]