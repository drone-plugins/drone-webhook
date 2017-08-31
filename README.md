# drone-webhook

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-webhook/status.svg)](http://beta.drone.io/drone-plugins/drone-webhook)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-webhook/coverage.svg)](https://aircover.co/drone-plugins/drone-webhook)
[![](https://badge.imagelayers.io/plugins/drone-webhook:latest.svg)](https://imagelayers.io/?images=plugins/drone-webhook:latest 'Get your own badge on imagelayers.io')

Drone plugin to send build status notifications via Webhook. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Binary

Build the binary using `./.drone.sh`:

### Usage

```
  docker run --rm \
    -e PLUGIN_URLS=https://hooks.somplace.com/endpoing/... \
    -e PLUGIN_HEADERS="HEADER1=value1" \
    -e PLUGIN_USERNAME=drone \
    -e PLUGIN_PASSWORD=password \
    -e DRONE_REPO_OWNER=octocat \
    -e DRONE_REPO_NAME=hello-world \
    -e DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
    -e DRONE_COMMIT_BRANCH=master \
    -e DRONE_COMMIT_AUTHOR=octocat \
    -e DRONE_BUILD_NUMBER=1 \
    -e DRONE_BUILD_STATUS=success \
    -e DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
    -e DRONE_TAG=1.0.0 \
    plugins/webhook
```






