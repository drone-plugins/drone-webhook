package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
)

func main() {
	var repo = drone.Repo{}
	var build = drone.Build{}
	var vargs = struct {
		Urls        []string          `json:"urls"`
		Headers     map[string]string `json:"header"`
		Method      string            `json:"method"`
		Template    string            `json:"template"`
		ContentType string            `json:"content_type"`
	}{}

	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.Parse()

	// data structure
	data := struct {
		Repo  drone.Repo  `json:"repo"`
		Build drone.Build `json:"build"`
	}{repo, build}

	// set default values
	if len(vargs.Method) == 0 {
		vargs.Method = "POST"
	}
	if len(vargs.ContentType) == 0 {
		vargs.ContentType = "application/json"
	}

	// creates the payload. by default the payload
	// is the build details in json format, but a custom
	// template may also be used.
	var buf bytes.Buffer
	if len(vargs.Template) == 0 {
		json.NewEncoder(&buf).Encode(&data)

	} else {
		t, err := template.New("_").Parse(vargs.Template)
		if err != nil {
			fmt.Printf("Error parsing content template. %s\n", err)
			os.Exit(1)
		}
		t.Execute(&buf, &data)
		if err != nil {
			fmt.Printf("Error executing content template. %s\n", err)
			os.Exit(1)
		}
	}

	// post payload to each url
	for _, rawurl := range vargs.Urls {

		uri, err := url.Parse(rawurl)
		if err != nil {
			fmt.Printf("Error parsing hook url. %s\n", err)
			os.Exit(1)
		}

		req, err := http.NewRequest(vargs.Method, uri.String(), &buf)
		if err != nil {
			fmt.Printf("Error creating http request. %s\n", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", vargs.ContentType)
		for key, value := range vargs.Headers {
			req.Header.Set(key, value)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error executing http request. %s\n", err)
			os.Exit(1)
		}
		resp.Body.Close()
	}
}
