FROM golang:1.23.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o run-app ./cmd

FROM debian:bookworm

WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY --from=builder /app/run-app .

COPY --from=builder /app/templates /root/templates

CMD ["./run-app"]
