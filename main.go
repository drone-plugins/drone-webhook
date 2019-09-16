package main

import (
	"os"

	"github.com/drone-plugins/drone-plugin-lib/pkg/urfave"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

const (
	// MethodFlag defines the method flag name
	MethodFlag = "method"

	// MethodEnvVar defines the method env var
	MethodEnvVar = "PLUGIN_METHOD"

	// UsernameFlag defines the username flag name
	UsernameFlag = "username"

	// UsernameEnvVar defines the username env var
	UsernameEnvVar = "PLUGIN_USERNAME,WEBHOOK_USERNAME"

	// PasswordFlag defines the password flag name
	PasswordFlag = "password"

	// PasswordEnvVar defines the password env var
	PasswordEnvVar = "PLUGIN_PASSWORD,WEBHOOK_PASSWORD"

	// ContentTypeFlag defines the content type flag name
	ContentTypeFlag = "content-type"

	// ContentTypeEnvVar defines the content type env var
	ContentTypeEnvVar = "PLUGIN_CONTENT_TYPE"

	// TemplateFlag defines the template flag name
	TemplateFlag = "template"

	// TemplateEnvVar defines the template env var
	TemplateEnvVar = "PLUGIN_TEMPLATE"

	// HeadersFlag defines the headers flag name
	HeadersFlag = "headers"

	// HeadersEnvVar defines the headers env var
	HeadersEnvVar = "PLUGIN_HEADERS"

	// URLsFlag defines the urls flag name
	URLsFlag = "urls"

	// URLsEnvVar defines the urls env var
	URLsEnvVar = "PLUGIN_URLS,PLUGIN_URL,WEBHOOK_URLS,WEBHOOK_URL"

	// ValidResponseCodesFlag defines the valid response codes flag name
	ValidResponseCodesFlag = "valid-response-codes"

	// ValidResponseCodesEnvVar defines the valid response codes env var
	ValidResponseCodesEnvVar = "PLUGIN_VALID_RESPONSE_CODES"

	// DebugFlag defines the debug flag name
	DebugFlag = "debug"

	// DebugEnvVar defines the debug env var
	DebugEnvVar = "PLUGIN_DEBUG"

	// SkipVerifyFlag defines the skip verify flag name
	SkipVerifyFlag = "skip-verify"

	// SkipVerifyEnvVar defines the skip verify env var
	SkipVerifyEnvVar = "PLUGIN_SKIP_VERIFY"
)

func main() {
	app := cli.NewApp()
	app.Name = "webhook plugin"
	app.Usage = "webhook plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   MethodFlag,
			Usage:  "webhook method",
			EnvVar: MethodEnvVar,
			Value:  "POST",
		},
		cli.StringFlag{
			Name:   UsernameFlag,
			Usage:  "username for basic auth",
			EnvVar: UsernameEnvVar,
		},
		cli.StringFlag{
			Name:   PasswordFlag,
			Usage:  "password for basic auth",
			EnvVar: PasswordEnvVar,
		},
		cli.StringFlag{
			Name:   ContentTypeFlag,
			Usage:  "content type",
			EnvVar: ContentTypeEnvVar,
			Value:  "application/json",
		},
		cli.StringFlag{
			Name:   TemplateFlag,
			Usage:  "custom template for webhook",
			EnvVar: TemplateEnvVar,
		},
		cli.StringSliceFlag{
			Name:   HeadersFlag,
			Usage:  "custom headers key map",
			EnvVar: HeadersEnvVar,
		},
		cli.StringSliceFlag{
			Name:   URLsFlag,
			Usage:  "list of urls to perform the action on",
			EnvVar: URLsEnvVar,
		},
		cli.IntSliceFlag{
			Name:   ValidResponseCodesFlag,
			Usage:  "list of valid http response codes",
			EnvVar: ValidResponseCodesEnvVar,
		},
		cli.BoolFlag{
			Name:   DebugFlag,
			Usage:  "enable debug information",
			EnvVar: DebugEnvVar,
		},
		cli.BoolFlag{
			Name:   SkipVerifyFlag,
			Usage:  "skip ssl verification",
			EnvVar: SkipVerifyEnvVar,
		},
	}

	flags := [][]cli.Flag{
		urfave.BuildFlags(),
		urfave.RepoFlags(),
		urfave.CommitFlags(),
		urfave.StageFlags(),
		urfave.StepFlags(),
		urfave.SemVerFlags(),
	}

	for _, flagz := range flags {
		app.Flags = append(
			app.Flags,
			flagz...,
		)
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Build:  urfave.BuildFromContext(c),
		Repo:   urfave.RepoFromContext(c),
		Commit: urfave.CommitFromContext(c),
		Stage:  urfave.StageFromContext(c),
		Step:   urfave.StepFromContext(c),
		SemVer: urfave.SemVerFromContext(c),

		Config: Config{
			Method:      c.String(MethodFlag),
			Username:    c.String(UsernameFlag),
			Password:    c.String(PasswordFlag),
			ContentType: c.String(ContentTypeFlag),
			Template:    c.String(TemplateFlag),
			Headers:     c.StringSlice(HeadersFlag),
			URLs:        c.StringSlice(URLsFlag),
			ValidCodes:  c.IntSlice(ValidResponseCodesFlag),
			Debug:       c.Bool(DebugFlag),
			SkipVerify:  c.Bool(SkipVerifyFlag),
		},
	}

	if plugin.Config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if len(plugin.Config.URLs) == 0 {
		log.Fatal("You must provide at least one url")
	}

	return plugin.Exec()
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
