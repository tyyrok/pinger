# Use the official Golang image as a builder
FROM golang:1.24.1-alpine3.20 as builder

# Set working directory inside the container
WORKDIR /app

# Copy go files
COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bot

FROM alpine:latest
COPY --from=builder /app/bot /app/bot

WORKDIR /app

