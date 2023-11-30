# Stage 1: Build the Go application
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todoapp main.go


# Stage 2: Create a minimal image to run the application
FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add curl

# Copy only the necessary files from the builder stage
COPY --from=builder /app/todoapp .


# Expose any necessary ports
EXPOSE 8080

# Command to run the application
CMD /app/todoapp
