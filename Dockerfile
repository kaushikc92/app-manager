FROM golang:1.12.5

WORKDIR $GOPATH/src/github.com/kaushikc92/app-manager

COPY ./app-manager ./

CMD ./app-manager
