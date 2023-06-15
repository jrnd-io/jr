FROM golang AS builder
MAINTAINER Ugo Landini <ugo@confluent.io>

ARG VERSION=0.3.0
ARG USER=$(id -u -n)
ARG TIME=$(date)

RUN useradd jr

WORKDIR /go/src/github.com/ugol/jr
COPY . .
RUN go get -u -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -tags static_all -v -ldflags="-X 'github.com/ugol/jr/pkg/cmd.Version=${VERSION}' -X 'github.com/ugol/jr/pkg/cmd.BuildUser=${USER}' -X 'github.com/ugol/jr/pkg/cmd.BuildTime=${TIME}'" -o build/jr jr.go

FROM registry.access.redhat.com/ubi9/ubi-micro
    
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/ugol/jr/templates/ /home/jr/.jr/templates/
COPY --from=builder /go/src/github.com/ugol/jr/config/ /home/jr/.jr/
COPY --from=builder /go/src/github.com/ugol/jr/pkg/producers/kafka/*.example /home/jr/.jr/kafka/
COPY --from=builder /go/src/github.com/ugol/jr/build/jr /bin

USER jr
