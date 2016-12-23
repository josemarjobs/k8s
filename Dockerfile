FROM ubuntu:trusty
MAINTAINER Josemar Magalhaes

ENV         PORT=3000
ENV         MONGO_SERVER_URL=mongodb
ENV         REDIS_SERVER_URL=redis

EXPOSE      $PORT
RUN         mkdir -p /var/server
ADD        ./dist/app /var/server
WORKDIR     /var/server
RUN         chmod +x app
ENTRYPOINT ./app
