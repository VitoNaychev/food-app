# Use an official Golang runtime as a parent image
FROM golang:1.21.3

# Set the working directory to /app
WORKDIR /app

# Copy the courier service at /app/courier-svc
COPY ../courier-svc /app/courier-svc

# Copy the authentication package at /app/auth
COPY ../auth /app/auth

# Copy the validation package at /app/validation
COPY ../validation /app/validation

# Copy the msgtypes package at /app/msgtypes
COPY ../msgtypes /app/msgtypes

# Copy the httperrors package at /app/httperrors
COPY ../httperrors /app/httperrors

# Copy the reqbuilder package at /app/reqbuilder
COPY ../reqbuilder /app/reqbuilder

# Copy go.mod and go.sum in /app
COPY ../go.mod /app
COPY ../go.sum /app

# Set the working directiry to /app/cmd
WORKDIR /app/courier-svc/cmd

# Build the Go application
RUN go build -o main

# Run the Go application
CMD ["./main"]