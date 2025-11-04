FROM golang:1.24-alpine AS builder

RUN mkdir -p /app
COPY . /app/
WORKDIR /app

# build
RUN go mod tidy && go build -o etl cmd/etl/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/etl etl
COPY --from=builder /app/build/config.yaml config/config.yaml
CMD ["./etl"]
