# multi-stage build
FROM golang:1.24-bookworm AS builder
WORKDIR /build
RUN apt-get update && apt-get install -y make && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN make build

FROM debian:bookworm-slim

# install tzdata for timezone support
RUN apt-get update && apt-get install -y tzdata && apt-get clean && rm -rf /var/lib/apt/lists/*

# by default, uses UTC timezone
# to change timezone, set TZ environment variable
# e.g. docker run -e TZ=America/New_York
# to verify the timezone inside the container, use the `date` command
ENV TZ=UTC
WORKDIR /app
COPY --from=builder /build/build/xcrawler /app/xcrawler

CMD [ "/app/xcrawler" ]