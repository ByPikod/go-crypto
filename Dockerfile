# Use the official Go image as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the container
COPY go.mod .
COPY go.sum .

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your Fiber application will listen on
EXPOSE 80

# Define the command to run your Fiber application when the container starts
CMD ["./main"]