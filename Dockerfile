FROM golang:1.20-alpine

WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main main.go

# Run the application
CMD ["/app/main"]
