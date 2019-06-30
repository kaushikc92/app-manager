FROM golang:1.12.5

WORKDIR $GOPATH/src/github.com/kaushikc92/app-manager

COPY ./app-manager ./

#COPY api.go ./
#COPY manager.go ./

#ENV GO111MODULE=on
#RUN go mod init

#RUN go get k8s.io/client-go@master

#RUN go build -o app-manager .

CMD ./app-manager
