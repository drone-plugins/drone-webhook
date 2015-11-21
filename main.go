package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
)

func main() {

	// plugin settings
	var repo = drone.Repo{}
	var build = drone.Build{}
	var vargs = Webhook{}

	// set plugin parameters
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)

	// parse the parameters
	if err := plugin.Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set default values
	if len(vargs.Method) == 0 {
		vargs.Method = "POST"
	}
	if len(vargs.ContentType) == 0 {
		vargs.ContentType = "application/json"
	}

	// data structure
	data := struct {
		Repo  drone.Repo  `json:"repo"`
		Build drone.Build `json:"build"`
	}{repo, build}

	// creates the payload. by default the payload
	// is the build details in json format, but a custom
	// template may also be used.
	var buf *bytes.Buffer
	var reqBytes []byte
	if len(vargs.Template) == 0 {
		json.NewEncoder(buf).Encode(&data)
		b := buf.Bytes()
		buf = bytes.NewBuffer(b)
		reqBytes = b
	} else {

		t, err := template.New("_").Parse(vargs.Template)
		if err != nil {
			fmt.Printf("Error parsing content template. %s\n", err)
			os.Exit(1)
		}

		t.Execute(buf, &data)
		if err != nil {
			fmt.Printf("Error executing content template. %s\n", err)
			os.Exit(1)
		}
		b := buf.Bytes()
		buf = bytes.NewBuffer(b)
		reqBytes = b
	}

	// build and execute a request for each url.
	// all auth, headers, method, template (payload),
	// and content_type values will be applied to
	// every webhook request.
	for i, rawurl := range vargs.Urls {

		uri, err := url.Parse(rawurl)
		if err != nil {
			fmt.Printf("Error parsing hook url. %s\n", err)
			os.Exit(1)
		}

		// vargs.Method defaults to POST, no need to check
		req, err := http.NewRequest(vargs.Method, uri.String(), buf)
		if err != nil {
			fmt.Printf("Error creating http request. %s\n", err)
			os.Exit(1)
		}

		// vargs.ContentType defaults to application/json, no need to check
		req.Header.Set("Content-Type", vargs.ContentType)
		for key, value := range vargs.Headers {
			req.Header.Set(key, value)
		}

		// set basic auth if a user or user and pass is provided
		if len(vargs.Auth.Username) > 0 {
			if len(vargs.Auth.Password) > 0 {
				req.SetBasicAuth(vargs.Auth.Username, vargs.Auth.Password)
			} else {
				req.SetBasicAuth(vargs.Auth.Username, "")
			}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error executing http request. %s\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if vargs.Verbose || os.Getenv("DEBUG") == "true" {

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// I do not think we need to os.Exit(1) if we are
				// unable to read a http response body.
				fmt.Printf("Error reading http response body. %s\n", err)
			}

			// scrub out basic auth pass
			/*
				if len(vargs.Auth.Password) > 0 {
					req.SetBasicAuth(vargs.Auth.Username, "XXXXX")
				}
			*/

			// print out
			fmt.Printf("Webhook URL %d\n  URL: %s\n  METHOD: %s\n  HEADERS: %s\n  BODY: %s\n  RESPONSE STATUS: %s\n  RESPONSE: %s\n", i+1, req.URL, req.Method, req.Header, string(reqBytes), resp.Status, string(body))
		}
	}
}
