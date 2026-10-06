# Build Stage
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage (Minimal Image)
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]