# drone-webhook
Drone plugin for sending Webhook notifications.

## Overview

This plugin is responsible for sending build notifications via Webhooks:

```sh
./drone-webhooks <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
      "urls": ["https://your.webhook/..."],
      "debug": true,
      "auth": {
        "username": "johnsmith",
        "password": "secretPass"
      },
      "headers": {
        "SomeHeader": "SomeHeaderValue"
      },
      "method": "POST",
      "template": "{\"git_branch\": \"{{ .Build.Branch }}\"}",
      "content_type": "application/json"
    }
}
EOF
```

## Docker

Build the Docker container. Note that we need to use the `-netgo` tag so that
the binary is built without a CGO dependency:

```sh
CGO_ENABLED=0 go build -a -tags netgo
docker build --rm=true -t plugins/drone-webhook .
```

Send a Webhook notification:

```sh
docker run -i plugins/drone-webhook <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
      "urls": ["https://your.webhook/..."],
      "debug": true,
      "auth": {
        "username": "johnsmith",
        "password": "secretPass"
      },
      "headers": {
        "SomeHeader": "SomeHeaderValue"
      },
      "method": "POST",
      "template": "{\"git_branch\": \"{{ .Build.Branch }}\"}",
      "content_type": "application/json"
    }
}
EOF
```
