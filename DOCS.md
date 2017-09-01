Use the Webhook plugin to notify services via Webhook when a build completes.
You will need to supply Drone with outgoing Webhook URLs.

You can override the default configuration with the following parameters:

* `urls` - JSON payloads are sent to each URL
* `method` - HTTP request method. Defaults to `POST`
* `headers` - HTTP request header map
* `username` - The username as a string for HTTP basic auth
* `password` - The password as a string for HTTP basic auth
* `skip_verify` - Skip verification of TLS certificates, defaults to `false`

## Example

The following is a sample configuration in your .drone.yml file:

```yaml
notify:
  webhook:
    urls:
      - https://your.webhook/...
      - https://your.other.webhook/...
    headers:
      - "Authorization=pa55word"
```

### Custom Body

In some cases you may want to submit a custom payload in the body of your hook.
For this usage the following additional parameters should be used:

* `template` - Handlebars template to create a custom payload body. See [docs](http://handlebarsjs.com/)
* `content_type` - HTTP request content type, defaults to `application/json`

Example configuration that generate a custom Yaml payload:

TBD

### Basic Authentication

> It is important to note that with HTTP Basic Authentication the provided
> username and password are not encrypted.

In some cases your webhook may need to authenticate with another service. You
can set the basic `Authentication` header with a username and password. For
these use cases we expose the following additional parameters:

* Sets the request's `Authorization` header to use HTTP Basic Authentication with the provided username and password below
  * `username` - The username as a string
  * `password` - The password as a string

Example configuration to include HTTP Basic Authentication:

```yaml
notify:
  webhook:
    method: POST
    username: myusername
    password: mypassword
    urls:
      - https://tower.example.com/...
```

### Debugging Webhooks

In some cases complicated webhooks may need debugging to ensure `urls`,
`template`, `auth` and more a properly configured. For these use cases we expose
the following `debug` parameter:

* `debug` - If `true` it will print out each URL request and response information

Example configuration to include the `debug` parameter:

```yaml
notify:
  webhook:
    debug: true
    method: POST
    username: myusername
    password: mypassword
    urls:
      - http://tower.example.com/api/v1/job_templates/44/launch/
      - http://tower.example.com/api/v1/job_templates/45/launch/
    content_type: application/json
```


