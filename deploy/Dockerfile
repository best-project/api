FROM alpine:3.8

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache curl

COPY . /root

EXPOSE 8080:8080

ENV PORT=8080
ENV DB_PORT=3306
ENV DB_USER=OMWVq72oqp
ENV DB_PASS=UcW1akjro5
ENV DB_NAME=OMWVq72oqp
ENV DB_HOST=remotemysql.com

ENTRYPOINT ["/root/api"]