FROM golang:1.11 as builder
ENV GOPATH /go
RUN go get -u github.com/golang/dep/cmd/dep \
    && mkdir -p $GOPATH/src/github.com/nstgt \
    && git clone https://github.com/nstgt/tapirjr.git $GOPATH/src/github.com/nstgt/tapirjr
WORKDIR $GOPATH/src/github.com/nstgt/tapirjr
RUN $GOPATH/bin/dep ensure \
    && go build -o $GOPATH/bin/tapirjr ./cmd/tapirjr

FROM alpine:3.8
WORKDIR /root
COPY --from=builder /go/bin/tapirjr /usr/local/bin/tapirjr
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
