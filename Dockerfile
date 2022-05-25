FROM ubuntu:18.04 AS builder
RUN mkdir -p /root/.ssh
RUN mkdir -p /etc/ssh/

USER root
RUN apt update && \
        #python \
        #py-pip \
        #groff \
        #less \
        #mailcap \
        #&& \s
    pip install --upgrade awscli==1.14.5 s3cmd==2.0.1 python-magic
   # apk -v --purge del py-pip && \
  #  rm /var/cache/apk/*
COPY id_rsa /root/.ssh/id_rsa
COPY --from=build ./lib/apm/appdynamics/libappdynamics.so /usr/local/lib/libappdynamics.so

#RUN apk add --update git openssh-client && \

RUN chmod 600 /root/.ssh/id_rsa && \
  eval $(ssh-agent) && \
  echo -e "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
  ssh-add /root/.ssh/id_rsa

# RUN apk update && apk add --no-cache git

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"

RUN DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates tzdata

WORKDIR /tokenizer
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -o /tokenizer/tokenizer
FROM scratch
COPY --from=builder /tokenizer /go/bin/tokenizer
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8083
ENTRYPOINT ["/go/bin/tokenizer/tokenizer", "start"]
