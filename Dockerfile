FROM debian:jessie

RUN apt-get update && apt-get install -y curl git

RUN curl -O https://storage.googleapis.com/golang/go1.3.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.3.2.linux-amd64.tar.gz

EXPOSE 8080
ENV PORT 8080

ENV GOPATH /opt/go
ADD . /opt/go/src/go-server-sent-events-example
WORKDIR /opt/go/src/go-server-sent-events-example
RUN /usr/local/go/bin/go get
RUN /usr/local/go/bin/go install

CMD /opt/go/bin/go-server-sent-events-example

