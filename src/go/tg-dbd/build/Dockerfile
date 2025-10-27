FROM golang:1.24-alpine AS builder

RUN mkdir -p /app
COPY . /app/
WORKDIR /app

# build
RUN go mod tidy && go build -o dbd cmd/dbd/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/dbd dbd
COPY --from=builder /app/build/config.yaml config/config.yaml
CMD ["./dbd"]
