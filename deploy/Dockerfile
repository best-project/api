FROM alpine:3.8

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache curl

COPY . /root

ENV DB_PORT=3306
ENV DB_USER=OMWVq72oqp
ENV DB_PASS=UcW1akjro5
ENV DB_NAME=OMWVq72oqp
ENV DB_HOST=remotemysql.com

ENV FB_APP_KEY=2595872383837249
ENV FB_APP_SECRET=c63fe1ea7aed77dc984daa198581bf76

ENTRYPOINT ["/root/api"]