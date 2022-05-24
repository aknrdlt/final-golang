FROM golang:latest

RUN mkdir /go/src/work
WORKDIR /go/src/work
ADD . /go/src/work

RUN go get -u github.com/gin-gonic/gin
CMD go run main.go