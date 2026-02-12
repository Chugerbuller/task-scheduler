FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web
COPY --from=builder /app/database ./database
COPY --from=builder /app/.env ./.env

EXPOSE 7540

CMD ["./scheduler"]