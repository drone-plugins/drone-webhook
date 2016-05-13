package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
	"github.com/drone/drone-go/template"
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

	system := drone.System{}
	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("system", &system)
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if vargs.Method == "" {
		vargs.Method = "POST"
	}

	if vargs.ContentType == "" {
		vargs.ContentType = "application/json"
	}

	// Creates the payload, by default the payload
	// is the build details in json format, but a custom
	// template may also be used.

	var buf bytes.Buffer

	if vargs.Template == "" {
		data := struct {
			System drone.System `json:"system"`
			Repo   drone.Repo   `json:"repo"`
			Build  drone.Build  `json:"build"`
		}{system, repo, build}

		if err := json.NewEncoder(&buf).Encode(&data); err != nil {
			fmt.Printf("Error: Failed to encode JSON payload. %s\n", err)
			os.Exit(1)
		}
	} else {
		err := template.Write(&buf, vargs.Template, &drone.Payload{
			Build:  &build,
			Repo:   &repo,
			System: &system,
		})

		if err != nil {
			fmt.Printf("Error: Failed to execute the content template. %s\n", err)
			os.Exit(1)
		}
	}

	// build and execute a request for each url.
	// all auth, headers, method, template (payload),
	// and content_type values will be applied to
	// every webhook request.

	for i, rawurl := range vargs.URLs {
		uri, err := url.Parse(rawurl)

		if err != nil {
			fmt.Printf("Error: Failed to parse the hook URL. %s\n", err)
			os.Exit(1)
		}

		b := buf.Bytes()
		r := bytes.NewReader(b)

		req, err := http.NewRequest(vargs.Method, uri.String(), r)

		if err != nil {
			fmt.Printf("Error: Failed to create the HTTP request. %s\n", err)
			os.Exit(1)
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
			fmt.Printf("Error: Failed to execute the HTTP request. %s\n", err)
			os.Exit(1)
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
}
