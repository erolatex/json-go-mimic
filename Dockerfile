# Build stage: Use the official Golang image to build the executable
FROM golang:latest as builder

# Set the working directory in the builder container
WORKDIR /build

# Copy the source code into the builder container
COPY src/*.go ./
COPY go.mod ./

# Download dependencies (if any)
RUN go mod download

# Build the executable with CGO disabled and optimization for Alpine Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o json-go-mimic .

# Run stage: Use a lightweight base image (Alpine)
FROM alpine:latest

# Install certificates for HTTPS communication
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the executable from the builder container
COPY --from=builder /build/json-go-mimic .

# Declare the port on which the server will run
EXPOSE 7732

# Command to run the executable
CMD ["./json-go-mimic"]

# Note: Configuration files and JSON mock files are expected to be mounted at runtime using Docker volumes.
# Example to run with mounted volumes:
# docker run -d -p 7732:7732 -v $(pwd)/configs:/app/configs -v $(pwd)/data:/app/data json-go-mimic