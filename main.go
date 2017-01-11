package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/template"
	"github.com/urfave/cli"
)

const (
	respFormat      = "Webhook %d\n  URL: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
	debugRespFormat = "Webhook %d\n  URL: %s\n  METHOD: %s\n  HEADERS: %s\n  REQUEST BODY: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
)

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone Webhook Plugin built from %s\n", buildCommit)

	app := cli.NewApp()
	app.Name = "webhook plugin"
	app.Usage = "webhook plugin"
	app.Action = run
	app.Version = fmt.Sprint(buildCommit)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "webhook debug",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.BoolFlag{
			Name:   "skip_verify",
			Usage:  "webhook skip verify",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.StringFlag{
			Name:   "method",
			Usage:  "webhook method",
			EnvVar: "PLUGIN_METHOD",
		},
		cli.StringSliceFlag{
			Name:   "urls",
			Usage:  "webhook urls",
			EnvVar: "PLUGIN_URLS",
		},
		cli.StringFlag{
			Name:   "content_type",
			Usage:  "webhook content type",
			EnvVar: "PLUGIN_CONTENT_TYPE",
		},
		cli.StringFlag{
			Name:   "auth.username",
			Usage:  "webhook auth username",
			EnvVar: "PLUGIN_AUTH_USERNAME",
		},
		cli.StringFlag{
			Name:   "auth.password",
			Usage:  "webhook auth password",
			EnvVar: "PLUGIN_AUTH_PASSWORD",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "webhook template",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.StringFlag{
			Name:  "repo.owner",
			Usage: "repo owner",
		},
		cli.StringFlag{
			Name:  "repo.name",
			Usage: "repo name",
		},
		cli.StringFlag{
			Name:  "repo.link",
			Usage: "repo link",
		},
		cli.StringFlag{
			Name:  "repo.avatar",
			Usage: "repo avatar",
		},
		cli.StringFlag{
			Name:  "repo.branch",
			Usage: "repo branch",
		},
		cli.StringFlag{
			Name:  "repo.clone",
			Usage: "repo clone",
		},
		cli.StringFlag{
			Name:  "commit.sha",
			Usage: "commit sha",
		},
		cli.StringFlag{
			Name:  "commit.ref",
			Usage: "commit ref",
		},
		cli.StringFlag{
			Name:  "commit.branch",
			Usage: "commit branch",
		},
		cli.StringFlag{
			Name:  "commit.link",
			Usage: "commit link",
		},
		cli.StringFlag{
			Name:  "commit.message",
			Usage: "commit message",
		},
		cli.StringFlag{
			Name:  "commit.author.name",
			Usage: "commit author name",
		},
		cli.StringFlag{
			Name:  "commit.author.email",
			Usage: "commit author email",
		},
		cli.StringFlag{
			Name:  "commit.author.avatar",
			Usage: "commit author avatar",
		},
		cli.IntFlag{
			Name:  "build.number",
			Usage: "build number",
		},
		cli.StringFlag{
			Name:  "build.event",
			Usage: "build event",
		},
		cli.StringFlag{
			Name:  "build.status",
			Usage: "build status",
		},
		cli.StringFlag{
			Name:  "build.link",
			Usage: "build link",
		},
		cli.Int64Flag{
			Name:  "build.created",
			Usage: "build created",
		},
		cli.Int64Flag{
			Name:  "build.started",
			Usage: "build started",
		},
		cli.Int64Flag{
			Name:  "build.finished",
			Usage: "build finished",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	logrus.Info(c.Generic("headers"))

	payload := &drone.Payload{
		Repo: &drone.Repo{
			Owner:  c.String("repo.owner"),
			Name:   c.String("repo.name"),
			Link:   c.String("repo.link"),
			Avatar: c.String("repo.avatar"),
			Branch: c.String("repo.branch"),
			Clone:  c.String("repo.clone"),
		},
		Build: &drone.Build{
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Link:     c.String("build.link"),
			Created:  c.Int64("build.created"),
			Started:  c.Int64("build.started"),
			Finished: c.Int64("build.finished"),
			Commit:   c.String("commit.sha"),
			Ref:      c.String("commit.ref"),
			Branch:   c.String("commit.branch"),
			Message:  c.String("commit.message"),
			Author:   c.String("commit.author.name"),
			Avatar:   c.String("commit.author.avatar"),
			Email:    c.String("commit.author.email"),
		},
	}

	vargs := Params{
		URLs:       c.StringSlice("urls"),
		SkipVerify: c.Bool("skip_verify"),
		Debug:      c.Bool("debug"),
		Auth: Auth{
			Password: c.String("auth.password"),
			Username: c.String("auth.username"),
		},
		Method:      c.String("method"),
		Template:    c.String("template"),
		ContentType: c.String("content_type"),
	}

	if vargs.Method == "" {
		vargs.Method = "POST"
	}

	if vargs.ContentType == "" {
		vargs.ContentType = "application/json"
	}

	var b []byte
	if vargs.Template == "" {
		buf, err := json.Marshal(&payload)
		if err != nil {
			return fmt.Errorf("Error: Failed to encode JSON payload. %s\n", err)
		}
		b = buf
	} else {
		msg, err := template.RenderTrim(vargs.Template, &payload)
		if err != nil {
			return fmt.Errorf("Error: Failed to execute the content template. %s\n", err)
		}
		b = []byte(msg)
	}

	// build and execute a request for each url.
	// all auth, headers, method, template (payload),
	// and content_type values will be applied to
	// every webhook request.

	for i, rawurl := range vargs.URLs {
		uri, err := url.Parse(rawurl)

		if err != nil {
			return fmt.Errorf("Error: Failed to parse the hook URL. %s\n", err)
		}

		r := bytes.NewReader(b)

		req, err := http.NewRequest(vargs.Method, uri.String(), r)

		if err != nil {
			return fmt.Errorf("Error: Failed to create the HTTP request. %s\n", err)
		}

		req.Header.Set("Content-Type", vargs.ContentType)

		for key, value := range vargs.Headers {
			req.Header.Set(key, value)
		}

		if vargs.Auth.Username != "" {
			req.SetBasicAuth(vargs.Auth.Username, vargs.Auth.Password)
		}

		client := http.DefaultClient
		if vargs.SkipVerify {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}
		resp, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("Error: Failed to execute the HTTP request. %s\n", err)
		}

		defer resp.Body.Close()

		if vargs.Debug || resp.StatusCode >= http.StatusBadRequest {
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				fmt.Printf("Error: Failed to read the HTTP response body. %s\n", err)
			}

			if vargs.Debug {
				fmt.Printf(
					debugRespFormat,
					i+1,
					req.URL,
					req.Method,
					req.Header,
					string(b),
					resp.Status,
					string(body),
				)
			} else {
				fmt.Printf(
					respFormat,
					i+1,
					req.URL,
					resp.Status,
					string(body),
				)
			}
		}
	}
	return nil
}
