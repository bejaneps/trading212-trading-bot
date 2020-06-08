FROM golang:latest

RUN mkdir /trading212; mkdir /trading212/logs

ADD . /trading212

WORKDIR /trading212

RUN ./build.sh

CMD ["/trading212/bin/./web", "--inifile"]