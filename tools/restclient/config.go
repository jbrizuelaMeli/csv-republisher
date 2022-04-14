package restclient

import "time"

type Config struct {
	ApiDomain        string                     `yaml:"api_domain"`
	TimeoutMillis    time.Duration              `yaml:"timeout"`
	ExternalApiCalls map[string]ExternalApiCall `yaml:"external_calls"`
}

type ExternalApiCall struct {
	ApiDomain string              `yaml:"domain"`
	Resources map[string]Resource `yaml:"resources"`
}

type Resource struct {
	RequestUri string        `yaml:"request_uri"`
	Auth       Authorization `yaml:"auth,omitempty"`
}

type Authorization struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
