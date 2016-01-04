package main

// Params represents the valid paramenter options for the webhook plugin.
type Params struct {
	Urls        []string          `json:"urls"`
	Debug       bool              `json:"debug"`
	Auth        Auth              `json:"auth"`
	Headers     map[string]string `json:"header"`
	Method      string            `json:"method"`
	Template    string            `json:"template"`
	ContentType string            `json:"content_type"`
}

// Auth represents a basic HTTP authentication username and password.
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
