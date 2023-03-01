FROM golang:alpine3.17

LABEL version="1.0"

LABEL maintainer="Nabcake Squad <administrator@astraea.jp.net>"

RUN apk add git gcc musl-dev

RUN git clone https://github.com/AdamH18/Scanlation-Discord-Bot.git

RUN cd ./Scanlation-Discord-Bot && CGO_ENABLED=1 GOOS=linux GOARCH=amd64\
    go build && mkdir /app && mv ./scanlation-discord-bot /app/scanlation-discord-bot && cp config-template.json /app/config.json

WORKDIR "/app"

CMD ["/app/scanlation-discord-bot"]
