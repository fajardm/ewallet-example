FROM golang:1.14

ENV GO111MODULE=on

WORKDIR /go/src/github.com/fajardm/ewallet-example
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["ewallet-example"]
