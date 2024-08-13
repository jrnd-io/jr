FROM jr-alpine as jr
FROM golang:1.20 as builder

# build with exec extension
RUN go install go.k6.io/xk6/cmd/xk6@latest && \
    xk6 build \
    --with github.com/grafana/xk6-exec@d9fdfce4ac85e6f0c42100794b0e51047f5dd0f6 \
    --output /k6


# build jr


# build from grafana
FROM grafana/k6:latest

USER jr-user
WORKDIR /home/jr

COPY --from=builder /k6 /usr/bin/k6
COPY --from=jr /etc/passwd /etc/passwd
COPY --from=jr /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=jr /home/jr/.jr .jr
COPY --from=jr /bin/jr /usr/bin/jr
