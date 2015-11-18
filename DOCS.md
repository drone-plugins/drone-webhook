Use the Webhook plugin to notify services via Webhook when a build completes.
You will need to supply Drone with outgoing Webhook URLs.

The following parameters are used to configure outgoing Webhooks:

* `urls` - JSON payloads are sent to each URL listed here
* `method` - HTTP request method. Defaults to `POST`
* `header` - HTTP request header map

The following is a sample Webhook configuration in your .drone.yml file:

```yaml
notify:
  webhook:
    urls:
      - https://your.webhook/...
      - https://your.other.webhook/...
    header:
      Authorization: pa55word
```

The following is an example Webhook payload (whitespace added):

```json
{
    "build": {
        "author": "johnsmith",
        "author_avatar": "",
        "author_email": "john.smith@gmail.com",
        "branch": "master",
        "commit": "9f2849d5",
        "created_at": 0,
        "enqueued_at": 0,
        "event": "",
        "finished_at": 1421029813,
        "link_url": "",
        "message": "Update the Readme",
        "number": 22,
        "ref": "",
        "refspec": "",
        "remote": "",
        "started_at": 1421029603,
        "status": "success",
        "timestamp": 0,
        "title": ""
    },
    "repo": {
        "allow_deploys": false,
        "allow_pr": false,
        "allow_push": false,
        "allow_tags": false,
        "avatar_url": "",
        "clone_url": "",
        "default_branch": "",
        "full_name": "foo/bar",
        "link_url": "",
        "name": "bar",
        "owner": "foo",
        "private": false,
        "timeout": 0,
        "trusted": false
    }
}
```

## Custom Body

In some cases you may want to submit a custom payload in the body of your hook. For the use case we expose the following additional parameters:

* `template` - Go template to create a custom payload body. See [docs](https://golang.org/pkg/text/template/)
* `content_type` - HTTP request content type

Example configuration that generate a custom Yaml payload:

```yaml
notify:
  webhook:
    urls:
      - https://your.webhook/...
      - https://your.other.webhook/...
    content_type: application/yaml
    template: >
      repo: {{.Repo.FullName}}
      build: {{.Build.Number}}
      commit: {{.Build.Commit}}
```