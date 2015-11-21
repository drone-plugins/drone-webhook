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

	fmt.Println(os.Getenv("DEBUG"))

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
	var buf bytes.Buffer
	if len(vargs.Template) == 0 {
		if err := json.NewEncoder(&buf).Encode(&data); err != nil {
			fmt.Printf("Error encoding content template. %s\n", err)
			os.Exit(1)
		}
	} else {

		t, err := template.New("_").Parse(vargs.Template)
		if err != nil {
			fmt.Printf("Error parsing content template. %s\n", err)
			os.Exit(1)
		}

		if err := t.Execute(&buf, &data); err != nil {
			fmt.Printf("Error executing content template. %s\n", err)
			os.Exit(1)
		}
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
		b := buf.Bytes()
		r := bytes.NewReader(b)
		req, err := http.NewRequest(vargs.Method, uri.String(), r)
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

		if vargs.Debug || os.Getenv("DEBUG") == "true" {

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// I do not think we need to os.Exit(1) if we are
				// unable to read a http response body.
				fmt.Printf("Error reading http response body. %s\n", err)
			}

			// debug print
			fmt.Printf("[debug] Webhook %d\n  URL: %s\n  METHOD: %s\n  HEADERS: %s\n  REQUEST BODY: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n", i+1, req.URL, req.Method, req.Header, string(b), resp.Status, string(body))

		} else if resp.StatusCode >= http.StatusBadRequest {

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// I do not think we need to os.Exit(1) if we are
				// unable to read a http response body.
				fmt.Printf("Error reading http response body. %s\n", err)
			}

			fmt.Printf("[info] Webhook %d\n  URL: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n", i+1, req.URL, resp.Status, string(body))
		}
	}
}
