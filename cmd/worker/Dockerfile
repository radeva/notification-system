FROM golang:1.24-alpine

# Copy the source code
COPY . .

RUN pwd && ls -la   

RUN go mod tidy

# Download all dependencies
RUN go mod download


# Build the application
RUN go build -o worker cmd/worker/main.go

# Command to run the executable
CMD ["./worker"] 