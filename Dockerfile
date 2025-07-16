FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./bin/marketplace ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/marketplace .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["/app/marketplace"]
