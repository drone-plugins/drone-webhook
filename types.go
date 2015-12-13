package main

type Params struct {
	Urls        []string          `json:"urls"`
	Debug       bool              `json:"debug"`
	Auth        Auth              `json:"auth"`
	Headers     map[string]string `json:"header"`
	Method      string            `json:"method"`
	Template    string            `json:"template"`
	ContentType string            `json:"content_type"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
