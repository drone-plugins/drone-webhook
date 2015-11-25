package main

// Webhook represents a webhook.
type Webhook struct {
	Urls        []string          `json:"urls"`
	Debug       bool              `json:"debug"`
	Auth        BasicAuth         `json:"auth"`
	Headers     map[string]string `json:"header"`
	Method      string            `json:"method"`
	Template    string            `json:"template"`
	ContentType string            `json:"content_type"`
}

// BasicAuth represents a HTTP basic authentication username and password.
type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
