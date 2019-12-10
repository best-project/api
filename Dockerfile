FROM alpine:3.8

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache curl

COPY . /root

RUN mkdir -p img/course

VOLUME img

ENTRYPOINT ["/root/api"]