# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang base image
FROM golang:1.23 as builder

# Add Maintainer Info
LABEL maintainer="Nam <namhq.1989@gmail.com>"

RUN mkdir -p /app

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Copy data to working dir
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/main ./cmd/*.go

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app .

# Command to run the executable
CMD ["./cmd/main"]
