FROM golang:1.16-alpine as build
WORKDIR /build
COPY src/* .
RUN go build -o go-away

FROM alpine:3.7 as app

RUN apk add --update supervisor redis 
COPY redis.conf /etc/redis.conf
COPY supervisor.conf /etc/supervisor.conf
COPY --from=build /build/go-away /app/go-away

EXPOSE 8080
CMD ["supervisord", "-c", "/etc/supervisor.conf"]

HEALTHCHECK CMD redis-cli ping