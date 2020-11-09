// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/drone/drone-template-lib/template"
	"github.com/urfave/cli/v2"
)

// Settings for the plugin.
type Settings struct {
	Method      string
	Username    string
	Password    string
	ContentType string
	Template    string
	Headers     cli.StringSlice
	URLs        cli.StringSlice
	ValidCodes  cli.StringSlice
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if len(p.settings.URLs.Value()) == 0 {
		return errors.New("must provide at least one webhook url")
	}

	for _, raw := range p.settings.URLs.Value() {
		_, err := url.Parse(raw)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	b, err := p.payload()

	if err != nil {
		return err
	}

	for i, raw := range p.settings.URLs.Value() {
		uri, _ := url.Parse(raw)

		req, err := http.NewRequest(p.settings.Method, uri.String(), bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("failed to create http request: %w", err)
		}

		req.Header.Set("Content-Type", p.settings.ContentType)
		for _, value := range p.settings.Headers.Value() {
			header := strings.SplitN(value, "=", 2)
			req.Header.Set(header[0], header[1])
		}

		if p.settings.Username != "" && p.settings.Password != "" {
			req.SetBasicAuth(p.settings.Username, p.settings.Password)
		}

		resp, err := p.network.Client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to execute http request: %w", err)
		}

		defer resp.Body.Close()

		if /*p.settings.Debug ||*/ resp.StatusCode >= http.StatusBadRequest {
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return fmt.Errorf("failed to parse http response: %w", err)
			}

			output, err := template.RenderTrim(result, webhook{
				Debug:    true, //p.settings.Debug,
				Number:   i,
				URL:      req.URL.String(),
				Method:   req.Method,
				Header:   req.Header,
				Status:   resp.Status,
				Request:  string(b),
				Response: string(body),
			})
			if err != nil {
				return fmt.Errorf("failed to parse debug template: %w", err)
			}

			fmt.Println(output)
		}

		validCodes := p.settings.ValidCodes.Value()
		if len(validCodes) > 0 && !intInSlice(validCodes, resp.StatusCode) {
			return fmt.Errorf("response of %d is not valid", resp.StatusCode)
		}
	}

	return nil
}

func (p *Plugin) payload() ([]byte, error) {
	if p.settings.Template == "" {
		res, err := json.Marshal(&p.pipeline)

		if err != nil {
			return []byte{}, fmt.Errorf("failed to generate json response: %w", err)
		}

		return res, nil
	}

	res, err := template.RenderTrim(p.settings.Template, p)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to parse response template: %w", err)
	}

	return []byte(res), nil
}
