FROM alpine:3.8

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache curl

COPY api /app/api

RUN mkdir -p images/

CMD ["/app/api"]