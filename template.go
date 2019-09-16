package main

import (
	"net/http"
)

// webhook defines the result payload.
type webhook struct {
	Number   int
	URL      string
	Method   string
	Header   http.Header
	Status   string
	Request  string
	Response string
	Debug    bool
}

// result defines the output template.
const result = `Webhook {{Number}}
  URL: {{URL}}{{#if Debug}}
  METHOD: {{Method}}
  HEADERS: {{Headers}}
  REQUEST: {{Request}}{{/if}}
  STATUS: {{Status}}
  RESPONSE: {{Response}}
`
