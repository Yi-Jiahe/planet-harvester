# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./src ./src

RUN go build -o ./bin/websocket-server ./src/websocket

CMD [ "./bin/websocket-server" ]
