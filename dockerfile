FROM alpine:latest

ARG PB_VERSION=0.14.0

RUN apk add --no-cache \
    unzip \
    openssh

ADD https://github.com/kasif-apps/backend/releases/download/v0.0.1/server.zip /tmp/pb.zip
RUN unzip /tmp/pb.zip -d /pb/

EXPOSE 8080

CMD ["/pb/server", "serve", "--http=0.0.0.0:8080"]