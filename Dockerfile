# Docker image for the Drone Webhook plugin
#
#     cd $GOPATH/src/github.com/drone-plugins/drone-webhook
#     make deps build docker

FROM alpine:3.3

RUN apk update && \
  apk add \
    ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-webhook /bin/
ENTRYPOINT ["/bin/drone-webhook"]
