FROM plugins/base:multiarch

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" \
  org.label-schema.name="Drone Webhook" \
  org.label-schema.vendor="Drone.IO Community" \
  org.label-schema.schema-version="1.0"

ADD release/linux/arm64/drone-webhook /bin/
ENTRYPOINT ["/bin/drone-webhook"]
