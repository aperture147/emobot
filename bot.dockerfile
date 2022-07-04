FROM golang:1.18-bullseye AS build-env

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    curl \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -gcflags="-l -l -l -l" -o bin/bot ./bin

FROM debian:bullseye-slim as run-env

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    curl \
    python \
    openssh-server \
    iproute2 \
 && rm -rf /var/lib/apt/lists/*

ADD ./heroku-exec.sh /app/.profile.d/heroku-exec.sh

RUN rm /bin/sh && ln -s /bin/bash /bin/sh

WORKDIR /run

COPY --from=build-env /build/bin/bot .

ENV BOT_TOKEN=placeholder \
    MONGO_URI=placeholder \
    BONSAI_URL=placeholder

CMD /run/bot