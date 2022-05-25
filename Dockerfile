FROM ubuntu:18.04 AS builder
RUN mkdir -p /root/.ssh
RUN mkdir -p /etc/ssh/

USER root

RUN apt-get update && \
    apt-get install -y \
        python3 \
        python3-pip \
        python3-setuptools \
        groff \
        less \
        golang-go \
    && pip3 install --upgrade pip \
    && apt-get clean

RUN pip3 --no-cache-dir install --upgrade awscli==1.14.5 s3cmd==2.0.1 python-magic

   # apk -v --purge del py-pip && \
  #  rm /var/cache/apk/*
COPY id_rsa /root/.ssh/id_rsa
COPY ./lib/apm/appdynamics/lib/libappdynamics.so /usr/local/lib/libappdynamics.so

RUN apt-get install -y openssh-client &&\
      apt-get install -y git
RUN chmod 600 /root/.ssh/id_rsa && \
  eval $(ssh-agent) && \
  echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
  ssh-add /root/.ssh/id_rsa

# RUN apk update && apk add --no-cache git

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"

RUN DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates tzdata

WORKDIR /tokenizer
COPY . .
RUN go get -d -v
# RUN CGO_ENABLED=1 GOOS=linux go build -o /tokenizer/tokenizer
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo

FROM scratch
COPY --from=builder /tokenizer /go/bin/tokenizer
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8083
ENTRYPOINT ["/go/bin/tokenizer/tokenizer", "start"]
