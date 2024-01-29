# Use an official Golang runtime as a parent image
FROM golang:1.21.3

# Set the working directory to /app
WORKDIR /app

# Copy the delivery service at /app/delivery-svc
COPY ../delivery-svc /app/delivery-svc

# Copy the authentication package at /app/auth
COPY ../auth /app/auth

# Copy the validation package at /app/validation
COPY ../validation /app/validation

# Copy the events package at /app/events
COPY ../events /app/events

# Copy the httperrors package at /app/httperrors
COPY ../httperrors /app/httperrors

# Copy the reqbuilder package at /app/reqbuilder
COPY ../reqbuilder /app/reqbuilder

# Copy the appenv package at /app/appenv
COPY ../appenv /app/appenv

# Copy the sm package at /app/sm
COPY ../sm /app/sm

# Copy the msgtypes package at /app/msgtypes
COPY ../msgtypes /app/msgtypes

# Copy the storeerror package at /app/storeerror
COPY ../storeerrors /app/storeerrors

# Copy the pgconfig package at /app/pgconfig
COPY ../pgconfig /app/pgconfig

# Copy go.mod and go.sum in /app
COPY ../go.mod /app
COPY ../go.sum /app

# Set the working directiry to /app/cmd
WORKDIR /app/delivery-svc/cmd

# Build the Go application
RUN go build -o main

# Run the Go application
CMD ["./main"]