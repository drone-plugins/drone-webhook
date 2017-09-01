FROM plugins/base:multiarch
MAINTAINER Drone.IO Community <drone-dev@googlegroups.com>

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/drone-plugins/drone-webhook.git"
LABEL org.label-schema.name="Drone Webhook"
LABEL org.label-schema.vendor="Drone.IO Community"
LABEL org.label-schema.schema-version="1.0"

ADD release/linux/amd64/drone-webhook /bin/
ENTRYPOINT ["/bin/drone-webhook"]
