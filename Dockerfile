FROM golang:1.22.4-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expose the port on which the application will run
EXPOSE 8000

# Run the application
CMD ["./main"]
