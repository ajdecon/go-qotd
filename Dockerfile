FROM google/debian:wheezy
MAINTAINER Adam DeConinck <ajdecon@ajdecon.org>

RUN apt-get update
RUN apt-get -y install golang git

ENV GOPATH /app
RUN mkdir /app
RUN go get github.com/ajdecon/go-qotd

EXPOSE 17
CMD ["/app/bin/go-qotd", "--file=/app/src/github.com/ajdecon/go-qotd/sample.data"]
