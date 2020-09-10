#конкретная версия
FROM golang:1.13 AS build

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GO111MODULE=on

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./main.go

#конкретная версия убунты
FROM ubuntu:20.04

EXPOSE 6000

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/main .

CMD ./main