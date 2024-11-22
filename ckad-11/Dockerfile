# Start with a lightweight Golang image
FROM golang:1.23.1-alpine

# Set the working directory
WORKDIR /app

# Copy source code and dependencies
COPY main.go .
COPY go.mod .
RUN go build -o app .

# Set environment variable
ENV SUN_CIPHER_ID=5b9c1065-e39d-4a43-a04a-e59bcea3e03f

# Run the application
CMD ["./app"]
