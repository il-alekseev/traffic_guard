# Build app
FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY .. .

RUN go mod download

RUN go build -o bisLog ./cmd/

# Small image
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates && mkdir config

COPY --from=builder /app/bisLog .
COPY --from=builder /app/docs .
#COPY ../config/config.yml ./config/config.yml

CMD ["./bisLog"]