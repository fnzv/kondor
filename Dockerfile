FROM golang:1.17.3

RUN mkdir -p /app

WORKDIR /app

ADD . /app

ENV GO111MODULE=on

RUN go mod init download

RUN go get ./...


RUN go build ./kondor.go

CMD ["./kondor"]
