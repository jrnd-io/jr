FROM golang AS builder
MAINTAINER Ugo Landini <ugo@confluent.io>

ARG VERSION=0.4.0
ARG GOVERSION=$(go version)
ARG USER=$(id -u -n)
ARG TIME=$(date)

RUN useradd jr

WORKDIR /go/src/github.com/jrnd-io/jr
COPY . .
RUN go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
RUN go generate pkg/generator/generate.go
RUN go get -u -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -tags static_all -v -ldflags="-X 'github.com/jrnd-io/jr/pkg/cmd.Version=${VERSION}' -X 'github.com/jrnd-io/jr/pkg/cmd.GoVersion=${GOVERSION}' -X 'github.com/jrnd-io/jr/pkg/cmd.BuildUser=${USER}' -X 'github.com/jrnd-io/jr/pkg/cmd.BuildTime=${TIME}'" -o build/jr jr.go

FROM registry.access.redhat.com/ubi9/ubi-micro
    
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/jrnd-io/jr/templates/ /home/jr/.jr/templates/
COPY --from=builder /go/src/github.com/jrnd-io/jr/config/ /home/jr/.jr/
COPY --from=builder /go/src/github.com/jrnd-io/jr/pkg/producers/kafka/*.example /home/jr/.jr/kafka/
COPY --from=builder /go/src/github.com/jrnd-io/jr/build/jr /bin

USER jr

ENV JR_SYSTEM_DIR=/home/jr/.jr