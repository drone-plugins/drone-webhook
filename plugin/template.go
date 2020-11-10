// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

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
