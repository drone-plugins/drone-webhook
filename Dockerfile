# Docker image for Drone's webhook notification plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-webhook .

FROM gliderlabs/alpine:3.2
RUN apk-install ca-certificates
ADD drone-webhook /bin/
ENTRYPOINT ["/bin/drone-webhook"]
