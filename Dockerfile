FROM golang:alpine

WORKDIR /app

COPY . /app

RUN curl -fsSL https://bun.sh/install | bash

RUN go generate
RUN go build -o filething main.go 

ENTRYPOINT ["./filething"]