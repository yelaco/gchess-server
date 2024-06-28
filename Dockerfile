# Build stage
FROM golang:1.22-alpine

WORKDIR /root/workspace/projects/go-chess-server

# Copy the source code into the container
COPY . .

# Copy go.mod and go.sum files first for better caching of dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

RUN go build -o go-chess-server ./cmd/game-server/server.go

# Expose the necessary ports
EXPOSE 80
EXPOSE 443
EXPOSE 7201
EXPOSE 7202

# Run the server
CMD ["./go-chess-server"]
