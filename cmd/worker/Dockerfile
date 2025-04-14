FROM golang:1.24-alpine

COPY go.mod ./
RUN go mod tidy

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o worker cmd/worker/main.go

# Command to run the executable
CMD ["./worker"] 