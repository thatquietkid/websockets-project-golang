# Use the official Go image on Alpine Linux for a small footprint.
FROM golang:1.24.4-alpine

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first.
# This leverages Docker's layer caching for faster builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Copy the frontend assets and TLS certificates into the working directory.
# Your Go application needs these to serve the UI and handle HTTPS.
COPY ./frontend ./frontend
COPY ./cert ./cert

# Build the Go application.
# CGO_ENABLED=0 creates a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .

# Expose port 8080, which is the port your Go app listens on.
EXPOSE 8080

# The command to run when the container starts.
# This directly executes your compiled Go application.
CMD ["/main"]
