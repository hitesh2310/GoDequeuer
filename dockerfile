# Stage 1: Build the binary
FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /go/bin/src/app

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/main ./cmd/main.go

# Stage 2: Copy the binary into a minimal Docker image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

RUN mkdir -p /app/logs/
# Copy the binary from the builder stage into the current working directory in the final image
COPY --from=builder  /go/bin/main /app/main
COPY --from=builder /go/bin/src/app/config/config.json /app/config/config.json
# Expose port 8081 to the outside world
EXPOSE 8081

# Command to run the executable
ENTRYPOINT ["./main"]

CMD [ "JOKER" ]
