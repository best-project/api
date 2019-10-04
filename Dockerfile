FROM alpine:3.8

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache curl

COPY ./api /root/api

ENTRYPOINT ["/root/api"]