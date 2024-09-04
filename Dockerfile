FROM golang:alpine

WORKDIR /app

COPY . /app

RUN apk add npm

RUN go generate
RUN go build main.go

ENTRYPOINT main.go