# Use the official Go image as a parent image
FROM golang:latest
LABEL org.opencontainers.image.source=https://github.com/ByPikod/Go-Crypto
LABEL org.opencontainers.image.authors="admin@yahyabatulu.com"

# Set the working directory inside the container
WORKDIR /app

# Git
RUN apt-get -y update
RUN apt-get -y install git
RUN git config --global --add safe.directory /app

# Copy the Go module files to the container
COPY go.mod .
COPY go.sum .

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Define the command to run your Fiber application when the container starts
CMD ["./main"]
