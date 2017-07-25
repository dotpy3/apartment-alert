FROM golang:latest

RUN go get github.com/kardianos/govendor
RUN govendor sync

RUN mkdir -p $GOPATH/src/github.com/dotpy3
RUN ln -s $PWD $GOPATH/src/github.com/dotpy3/apartment-alert
RUN pushd $GOPATH/src/github.com/dotpy3/apartment-alert

RUN go build -o /apartment-alert
RUN chmod +x /apartment-alert

ENTRYPOINT /apartment-alert
