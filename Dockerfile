FROM golang:1.20-alpine AS builder
MAINTAINER Ugo Landini <ugo@confluent.io>

ARG VERSION=0.1.0
ARG USER=$(id -u -n)
ARG TIME=$(date)

RUN apk update \
    && apk add --no-cache git \
    && apk add --no-cache ca-certificates \
    && apk add --update gcc musl-dev \
    && apk add --update librdkafka \
    && update-ca-certificates

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/jr" \
    --shell "/bin/sh" \
    --uid "100001" \
    "jr-user"

WORKDIR /go/src/github.com/ugol/jr
COPY . .

RUN go get -u -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -v -ldflags="-X 'jr/cmd.Version=${VERSION}' -X 'jr/cmd.BuildUser=${USER}'" -o build/jr jr.go

FROM alpine
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/ugol/jr/templates/ /home/jr/.jr/templates/
COPY --from=builder /go/src/github.com/ugol/jr/kafka/ /home/jr/.jr/kafka/
COPY --from=builder /go/src/github.com/ugol/jr/build/jr /bin

USER jr-user:jr-user
