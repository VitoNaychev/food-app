# Use an official Golang runtime as a parent image
FROM golang:1.21.3

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Set the working directiry to /app/cmd
WORKDIR /app/cmd

# Build the Go application
RUN go build -o main

# Run the Go application
CMD ["./main"]