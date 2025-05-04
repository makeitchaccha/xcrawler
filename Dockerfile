# multi-stage build
FROM golang:1.24 AS builder
WORKDIR /build
RUN apt-get update && apt-get install -y make && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN make build

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/build/xcrawler /app/xcrawler

CMD [ "/app/xcrawler" ]