FROM golang:1.24-alpine

# Copy the source code
COPY . .

RUN go mod tidy

# Download all dependencies
RUN go mod download

# Build the application
RUN go build -o api cmd/api/main.go

# Expose port 8081
EXPOSE 8081

# Command to run the executable
CMD ["./api"] 