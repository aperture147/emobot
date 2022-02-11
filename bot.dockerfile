FROM golang:1.17.6-alpine3.15 AS build-env

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -gcflags="-l=5" -o bot ./bot

FROM alpine:3.15.0

WORKDIR /run

COPY --from=build-env /build/bot/bot .

ENV BOT_TOKEN=placeholder \
    GUILD_ID=placeholder \
    MONGO_URI=placeholder

ENTRYPOINT bot