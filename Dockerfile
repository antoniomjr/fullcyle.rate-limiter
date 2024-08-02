# Use the official Go image as a parent image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Fetch app dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

# Copy the wait-for-it script
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Command to run the executable
CMD ["/wait-for-it.sh", "redis", "--", "./main"]