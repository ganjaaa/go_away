FROM golang:1.16-alpine
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY *.go ./

RUN go build -o /docker-go-away

EXPOSE 8080

CMD [ "/docker-go-away" ]