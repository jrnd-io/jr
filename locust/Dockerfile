FROM jr-alpine as jr
FROM locustio/locust as locust

USER jr-user
WORKDIR /home/jr


# add jr
COPY --from=jr /etc/passwd /etc/passwd
COPY --from=jr /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=jr /home/jr/.jr .jr
COPY --from=jr /bin/jr /usr/bin/jr
