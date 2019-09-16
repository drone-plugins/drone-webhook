package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/drone-plugins/drone-plugin-lib/pkg/plugin"
	"github.com/drone/drone-template-lib/template"
	"github.com/jackspirou/syscerts"
	log "github.com/sirupsen/logrus"
)

type (
	// Config provides the pugin configuration.
	Config struct {
		Method      string
		Username    string
		Password    string
		ContentType string
		Template    string
		Headers     []string
		URLs        []string
		ValidCodes  []int
		Debug       bool
		SkipVerify  bool
	}

	// Plugin provides all required attributes.
	Plugin struct {
		Build  plugin.Build
		Repo   plugin.Repo
		Commit plugin.Commit
		Stage  plugin.Stage
		Step   plugin.Step
		SemVer plugin.SemVer
		Config Config
	}
)

// Exec provides the concrete plugin handler.
func (p Plugin) Exec() error {
	b, err := p.payload()

	if err != nil {
		return err
	}

	for i, raw := range p.Config.URLs {
		uri, err := url.Parse(raw)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to parse hook URL")

			return err
		}

		req, err := http.NewRequest(p.Config.Method, uri.String(), bytes.NewReader(b))

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to create HTTP request")

			return err
		}

		req.Header.Set("Content-Type", p.Config.ContentType)

		for _, value := range p.Config.Headers {
			header := strings.SplitN(value, "=", 2)
			req.Header.Set(header[0], header[1])
		}

		if p.Config.Username != "" && p.Config.Password != "" {
			req.SetBasicAuth(p.Config.Username, p.Config.Password)
		}

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					RootCAs:            syscerts.SystemRootsPool(),
					InsecureSkipVerify: p.Config.SkipVerify,
				},
			},
		}

		resp, err := client.Do(req)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to execute HTTP request")

			return err
		}

		defer resp.Body.Close()

		if p.Config.Debug || resp.StatusCode >= http.StatusBadRequest {
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Failed to parse HTTP response")
			}

			output, err := template.RenderTrim(result, webhook{
				Debug:    p.Config.Debug,
				Number:   i,
				URL:      req.URL.String(),
				Method:   req.Method,
				Header:   req.Header,
				Status:   resp.Status,
				Request:  string(b),
				Response: string(body),
			})

			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Failed to parse debug template")

				return err
			}

			fmt.Println(output)
		}

		if len(p.Config.ValidCodes) > 0 && !intInSlice(p.Config.ValidCodes, resp.StatusCode) {
			log.WithFields(log.Fields{
				"code": resp.StatusCode,
			}).Error("Valid response code not found")

			return fmt.Errorf("valid response code not found")
		}
	}

	return nil
}

func (p Plugin) payload() ([]byte, error) {
	if p.Config.Template == "" {
		res, err := json.Marshal(&p)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to generate JSON response")

			return []byte{}, err
		}

		return res, nil
	}

	res, err := template.RenderTrim(p.Config.Template, p)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to parse response template")

		return []byte{}, err
	}

	return []byte(res), nil
}
