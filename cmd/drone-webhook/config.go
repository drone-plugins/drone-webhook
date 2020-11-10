// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/drone-plugins/drone-webhook/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	// Replace below with all the flags required for the plugin.
	// Use Destination within the cli.Flags to populate settings
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "method",
			Usage:       "webhook method",
			Value:       "POST",
			EnvVars:     []string{"PLUGIN_METHOD"},
			Destination: &settings.Method,
		},
		&cli.StringFlag{
			Name:        "username",
			Usage:       "username for basic auth",
			EnvVars:     []string{"PLUGIN_USERNAME", "WEBHOOK_USERNAME"},
			Destination: &settings.Username,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "password for basic auth",
			EnvVars:     []string{"PLUGIN_PASSWORD", "WEBHOOK_PASSWORD"},
			Destination: &settings.Password,
		},
		&cli.StringFlag{
			Name:        "content-type",
			Usage:       "content type",
			Value:       "application/json",
			EnvVars:     []string{"PLUGIN_CONTENT_TYPE"},
			Destination: &settings.ContentType,
		},
		&cli.StringFlag{
			Name:        "template",
			Usage:       "custom template for webhook",
			EnvVars:     []string{"PLUGIN_TEMPLATE"},
			Destination: &settings.Template,
		},
		&cli.StringSliceFlag{
			Name:        "headers",
			Usage:       "custom headers key map",
			EnvVars:     []string{"PLUGIN_HEADERS"},
			Destination: &settings.Headers,
		},
		&cli.StringSliceFlag{
			Name:        "urls",
			Usage:       "list of urls to perform the action on",
			EnvVars:     []string{"PLUGIN_URLS", "PLUGIN_URL", "WEBHOOK_URLS", "WEBHOOK_URL"},
			Destination: &settings.URLs,
		},
		// Should be an IntSliceFlag but that doesn't have a Destination field
		&cli.StringSliceFlag{
			Name:        "valid-response-codes",
			Usage:       "list of valid http response codes",
			EnvVars:     []string{"PLUGIN_VALID_RESPONSE_CODES"},
			Destination: &settings.ValidCodes,
		},
	}
}
