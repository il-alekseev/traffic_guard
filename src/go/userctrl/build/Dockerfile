FROM golang:1.24-alpine AS builder

LABEL stage=gobuilder
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o usercontrol cmd/usercontrol/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/usercontrol /app/usercontrol
COPY --from=builder /app/docs/ /app/docs/
COPY --from=builder /app/config/config.yml /app/config/config.yml
COPY --from=builder /app/pkg/clixonclient/nacm/ /app/pkg/clixonclient/nacm/

CMD ["./usercontrol"]