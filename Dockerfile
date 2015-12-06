FROM golang:latest

MAINTAINER Dmitry Kravtsov <idkravitz@gmail.com>

RUN wget https://storage.googleapis.com/kubernetes-release/release/v1.1.2/bin/linux/amd64/kubectl && chmod +x ./kubectl

ADD . .

RUN go install kube-reload
EXPOSE 80

ENTRYPOINT ["./bin/kube-reload"]