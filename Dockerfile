FROM golang:1.12.5

WORKDIR $GOPATH/src/github.com/kaushikc92/app-manager

ENV GO111MODULE=on
RUN go build

CMD ./app-manager
