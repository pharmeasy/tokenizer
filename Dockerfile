FROM golang:1.14.4-alpine3.11 AS builder
RUN mkdir -p /root/.ssh
RUN mkdir -p /etc/ssh/
RUN apk -v --update add \
        python \
        py-pip \
        groff \
        less \
        mailcap \
        && \
    pip install --upgrade awscli==1.14.5 s3cmd==2.0.1 python-magic && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/*
COPY id_rsa /root/.ssh/id_rsa
COPY ./lib/apm/appdynamics/lib/libappdynamics.so /usr/local/lib/libappdynamics.so
RUN apk add --update git openssh-client && \
chmod 600 /root/.ssh/id_rsa && \
  eval $(ssh-agent) && \
  echo -e "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
  ssh-add /root/.ssh/id_rsa
RUN apk add build-base
RUN apk update && apk add --no-cache git

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"

RUN apk --no-cache add ca-certificates

WORKDIR /tokenizer
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -o /tokenizer/tokenizer



FROM ubuntu:18.04
COPY --from=builder /tokenizer /go/bin/tokenizer
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8083
ENTRYPOINT ["/go/bin/tokenizer/tokenizer", "start"]

