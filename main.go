package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "webhook plugin"
	app.Usage = "webhook plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "method",
			Usage:  "webhook method",
			EnvVar: "PLUGIN_METHOD",
			Value:  "POST",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "username for basic auth",
			EnvVar: "PLUGIN_USERNAME,WEBHOOK_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "password for basic auth",
			EnvVar: "PLUGIN_PASSWORD,WEBHOOK_PASSWORD",
		},
		cli.StringFlag{
			Name:   "token-value",
			Usage:  "token value",
			EnvVar: "PLUGIN_TOKEN_VALUE,WEBHOOK_TOKEN_VALUE",
		},
		cli.StringFlag{
			Name:   "token-type",
			Usage:  "type of token",
			EnvVar: "PLUGIN_TOKEN_TYPE,WEBHOOK_TOKEN_TYPE",
			Value:  "Bearer",
		},
		cli.StringFlag{
			Name:   "content-type",
			Usage:  "content type",
			EnvVar: "PLUGIN_CONTENT_TYPE",
			Value:  "application/json",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "custom template for webhook",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.StringSliceFlag{
			Name:   "headers",
			Usage:  "custom headers key map",
			EnvVar: "PLUGIN_HEADERS",
		},
		cli.StringSliceFlag{
			Name:   "urls",
			Usage:  "list of urls to perform the action on",
			EnvVar: "PLUGIN_URLS,WEBHOOK_URLS",
		},
		cli.IntSliceFlag{
			Name:   "valid-response-codes",
			Usage:  "list of valid http response codes",
			EnvVar: "PLUGIN_VALID_RESPONSE_CODES",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "enable debug information",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.BoolFlag{
			Name:   "skip-verify",
			Usage:  "skip ssl verification",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.StringFlag{
			Name:   "signature-header",
			Usage:  "header name to use in request",
			EnvVar: "PLUGIN_SIGNATURE_HEADER,WEBHOOK_SIGNATURE_HEADER",
			Value:  "X-Drone-Signature",
		},
		cli.StringFlag{
			Name:   "signature-secret",
			Usage:  "secret to generate signature",
			EnvVar: "PLUGIN_SIGNATURE_SECRET,WEBHOOK_SIGNATURE_SECRET",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "build.deployTo",
			Usage:  "environment deployed to",
			EnvVar: "DRONE_DEPLOY_TO",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
		},
		cli.StringFlag{
			Name:   "stage.status",
			Usage:  "stage status",
			EnvVar: "DRONE_STAGE_STATUS",
		},
		cli.StringFlag{
			Name:   "stage.name",
			Usage:  "stage name",
			EnvVar: "DRONE_STAGE_NAME",
		},
		cli.StringFlag{
			Name:   "stage.type",
			Usage:  "stage type",
			EnvVar: "DRONE_STAGE_TYPE",
		},
		cli.StringFlag{
			Name:   "stage.kind",
			Usage:  "stage kind",
			EnvVar: "DRONE_STAGE_KIND",
		},
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env")
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Commit:   c.String("commit.sha"),
			Ref:      c.String("commit.ref"),
			Branch:   c.String("commit.branch"),
			Author:   c.String("commit.author"),
			Message:  c.String("commit.message"),
			Link:     c.String("build.link"),
			DeployTo: c.String("build.deployTo"),
			Started:  c.Int64("build.started"),
			Created:  c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			Method:          c.String("method"),
			Username:        c.String("username"),
			Password:        c.String("password"),
			TokenValue:      c.String("token-value"),
			TokenType:       c.String("token-type"),
			ContentType:     c.String("content-type"),
			Template:        c.String("template"),
			Headers:         c.StringSlice("headers"),
			URLs:            c.StringSlice("urls"),
			ValidCodes:      c.IntSlice("valid-response-codes"),
			Debug:           c.Bool("debug"),
			SkipVerify:      c.Bool("skip-verify"),
			SignatureHeader: c.String("signature-header"),
			SignatureSecret: c.String("signature-secret"),
		},
		Stage: Stage{
			Type:   c.String("stage.type"),
			Name:   c.String("stage.name"),
			Status: c.String("stage.status"),
			Kind:   c.String("stage.kind"),
		},
	}

	return plugin.Exec()
}
