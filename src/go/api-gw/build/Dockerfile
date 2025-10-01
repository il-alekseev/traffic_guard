FROM golang:1.24-alpine AS builder

RUN apk update --no-cache && apk add --no-cache gcc musl-dev

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

FROM alpine:latest

RUN apk add --no-cache libc6-compat
RUN apk --no-cache add ca-certificates curl krb5-libs krb5

WORKDIR /app
# copy executable
COPY --from=builder /app/server ./server
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/config ./config
COPY --from=builder /app/data ./data

EXPOSE 8001

CMD ["/app/server"]