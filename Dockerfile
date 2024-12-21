# Use the official Golang image as the base image
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application code to the container
COPY . .

# Install any required tools, if necessary (e.g., fsnotify for file watching)
RUN go install github.com/fsnotify/fsnotify@latest
RUN go install github.com/fatih/color@latest
RUN go insatll github.com/spf13/cobra@latest
# Command to run the application
CMD ["go", "run", "main.go"]
