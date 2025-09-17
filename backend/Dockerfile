# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /social-network ./server.go

# Stage 2: Create the final, smaller image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /social-network .
COPY --from=builder /app/imageUpload ./imageUpload
COPY --from=builder /app/pkg/db/migration ./pkg/db/migration

EXPOSE 8081

CMD ["./social-network"]