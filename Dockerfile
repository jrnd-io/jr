FROM alpine:latest AS build
MAINTAINER Gianluca Natali <gnatali@confluent.io>
RUN apk add --no-cache --update go gcc g++
WORKDIR .
ENV GOPATH /app
#COPY src /app/src
COPY . /app/src
#RUN go get jr 
RUN /bin/sh -c '/app/src/make_install.sh'

#server is name of our application
RUN CGO_ENABLED=1 GOOS=linux go install -a server

FROM alpine:latest
WORKDIR /app
RUN cd /app
COPY --from=build /app/src/build/jr /app/bin/jr
CMD ["bin/jr"]

ENTRYPOINT ["/bin/sh"]