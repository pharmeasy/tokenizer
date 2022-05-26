FROM golang:1.13.0 AS builder
RUN mkdir -p /root/.ssh
RUN mkdir -p /etc/ssh/
RUN apt-get update && \
    apt-get install -y \
        python3 \
        python3-pip \
        python3-setuptools \
    && pip3 install --upgrade pip \
    && apt-get clean

RUN pip3 --no-cache-dir install --upgrade awscli==1.14.5 s3cmd==2.0.1 python-magic
COPY id_rsa /root/.ssh/id_rsa
COPY ./lib/apm/appdynamics/lib/libappdynamics.so /usr/local/lib/libappdynamics.so

RUN apt-get install -y openssh-client &&\
      apt-get install -y git
RUN chmod 600 /root/.ssh/id_rsa && \
  eval $(ssh-agent) && \
  echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
  ssh-add /root/.ssh/id_rsa

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
RUN apt-get install -y ca-certificates

WORKDIR /tokenizer
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -o /tokenizer/tokenizer



FROM scratch
COPY --from=builder /tokenizer /go/bin/tokenizer
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder usr/local/lib/libappdynamics.so /usr/local/lib/libappdynamics.so

ENV LD_LIBRARY_PATH "/usr/local/lib"

EXPOSE 8083
ENTRYPOINT ["/go/bin/tokenizer/tokenizer", "start"]

