FROM ubuntu:trusty
MAINTAINER Josemar Magalhaes

ENV         PORT=5000
ENV         MONGO_SERVER_URL=mongodb
ENV         REDIS_SERVER_URL=redis

EXPOSE      $PORT
RUN         mkdir -p /var/server
ADD        ./dist/linux-amd64/app /var/server
WORKDIR     /var/server
ENTRYPOINT ./app
