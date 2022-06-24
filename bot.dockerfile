FROM golang:1.18.3-alpine3.16 AS build-env

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -gcflags="-l -l -l -l" -o bot ./bot

FROM alpine:3.16 as run-env

WORKDIR /run

COPY --from=build-env /build/bot/bot .

ENV BOT_TOKEN=placeholder \
    GUILD_ID=placeholder \
    MONGO_URI=placeholder

ENTRYPOINT /run/bot