FROM alpine:3.7

RUN apk add --update bash curl

COPY ./bin/linux-amd64/linkserver /usr/local/bin/linkserver
WORKDIR /
ENTRYPOINT ["/usr/local/bin/linkserver"]
