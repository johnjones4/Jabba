FROM golang:1.16

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/server

RUN go get ./...
RUN go build .

CMD ["/go/src/app/server/main"]
