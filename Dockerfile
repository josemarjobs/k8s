FROM ubuntu:trusty
MAINTAINER Josemar Magalhaes

ENV         PORT=3000
ENV         REDIS_SERVER_URL=redis

EXPOSE      $PORT
RUN         mkdir -p /var/server
ADD        ./dist/linux-amd64/app /var/server
WORKDIR     /var/server
ENTRYPOINT ./app

# gox \
#   -os="linux" \
#   -arch="amd64" \
#   -output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" app