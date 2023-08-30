FROM golang:alpine3.17 AS builder

LABEL version="1.1"
LABEL maintainer="Nabcake Squad <administrator@astraea.jp.net>"

RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o scanlation-discord-bot .

FROM alpine:3.17

RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/scanlation-discord-bot .

CMD ["./scanlation-discord-bot"]
