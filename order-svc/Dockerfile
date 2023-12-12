# Use an official Golang runtime as a parent image
FROM golang:1.21.3

# Set the working directory to /app
WORKDIR /app

# Copy the order service at /app/order-svc
COPY ../order-svc /app/order-svc

# Copy the validation package at /app/validation
COPY ../validation /app/validation

# Copy the msgtypes package at /app/msgtypes
COPY ../msgtypes /app/msgtypes

# Copy the httperrors package at /app/httperrors
COPY ../httperrors /app/httperrors

# Copy go.mod and go.sum in /app
COPY ../go.mod /app
COPY ../go.sum /app


# Set the working directiry to /app/cmd
WORKDIR /app/order-svc/cmd

# Build the Go application
RUN go build -o main

# # Run the Go application
CMD ["./main"]