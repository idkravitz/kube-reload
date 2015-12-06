FROM golang:latest

MAINTAINER Dmitry Kravtsov <idkravitz@gmail.com>

ADD . .

RUN go install kube-reload
EXPOSE 80

ENTRYPOINT ["./bin/kube-reload"]