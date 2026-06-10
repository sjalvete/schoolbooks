FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o schoolbooks ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=builder /app/schoolbooks .
RUN mkdir -p /app/data
VOLUME ["/app/data"]
EXPOSE 8080
CMD ["./schoolbooks"]
