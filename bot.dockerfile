FROM golang:1.18-alpine3.16 AS build-env

RUN apk --no-cache add curl

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -gcflags="-l -l -l -l" -o bin/bot ./bin

FROM alpine:3.16 as run-env

RUN rm /bin/sh \
 && ln -s /bin/bash /bin/sh \
 && mkdir -p /app/.profile.d/ \
 && printf '#!/usr/bin/env bash\n\nset +o posix\n\n[ -z "$SSH_CLIENT" ] && source <(curl --fail --retry 7 -sSL "$HEROKU_EXEC_URL")\n' > /app/.profile.d/heroku-exec.sh \
 && chmod +x /app/.profile.d/heroku-exec.sh \
 && ln -s /usr/bin/python3 /usr/bin/python

WORKDIR /run

COPY --from=build-env /build/bin/bot .

ENV BOT_TOKEN=placeholder \
    MONGO_URI=placeholder

ENTRYPOINT /run/bot