FROM golang:1.24-alpine AS builder

RUN mkdir -p /app
COPY . /app/
WORKDIR /app

# build
RUN go mod tidy && go build -o an cmd/an/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/an an
COPY --from=builder /app/build/config.yaml config/config.yaml
CMD ["./an"]
