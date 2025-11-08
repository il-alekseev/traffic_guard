# syntax=docker/dockerfile:1.7
FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/webscraper ./cmd/scraper

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

COPY --from=builder /tmp/webscraper /usr/local/bin/webscraper
COPY --from=builder /app/config ./config
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/data ./data

ENV SCRAPER_CONFIG="config/config.yaml"

EXPOSE 8010

ENTRYPOINT ["/usr/local/bin/webscraper", "--config", "config/config.yaml"]
