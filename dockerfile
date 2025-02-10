# Step 1: Use the official Golang image as the base image
FROM golang:1.22.5-alpine

# Step 2: Set the working directory inside the container
WORKDIR /Forum

# Step 3: Copy the go.mod and go.sum files first to leverage Docker's caching mechanism
COPY go.mod go.sum ./

# Step 4: Run `go mod tidy` to download dependencies
RUN go mod tidy

# Step 5: Copy the rest of your application source code into the container
COPY . .

# Step 6: Build your Go application
RUN go build -o main ./cmd/main.go

# Step 7: Expose port 8080 (or your application's port)
EXPOSE 8080

# Step 8: Define the default command to run your application
CMD ["./main"]
