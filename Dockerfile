# The FROM variable will be dynamically fetched via jenkins on staging and prod env
# Refer to the file tokenizer-base-DockerFile for building the project via docker on local env.
COPY id_rsa /root/.ssh/id_rsa

RUN chmod 600 /root/.ssh/id_rsa && \
  eval $(ssh-agent) && \
  echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
  ssh-add /root/.ssh/id_rsa

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"

WORKDIR /tokenizer
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -o /tokenizer/tokenizer



FROM ubuntu:18.04
COPY --from=builder /tokenizer /go/bin/tokenizer
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8083
ENTRYPOINT ["/go/bin/tokenizer/tokenizer", "start"]
