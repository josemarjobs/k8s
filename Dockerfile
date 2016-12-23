FROM ubuntu:trusty
MAINTAINER Josemar Magalhaes

ENV         PORT=3000
ENV         MONGO_SERVER_URL=mongodb

EXPOSE      $PORT
RUN         mkdir -p /var/server
ADD        ./dist/app /var/server
WORKDIR     /var/server
entrypoint  ./app
