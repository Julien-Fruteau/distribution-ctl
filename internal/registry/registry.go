package registry

import (
	"github.com/julien-fruteau/go/distctl/internal/env"
	httpUtils "github.com/julien-fruteau/go/distctl/internal/http"
)

var (
	host     = env.GetEnvOrDefault("REGISTRY_HOST", "localhost")
	scheme   = env.GetEnvOrDefault("REGISTRY_SCHEME", "http")
	username = env.GetEnvOrDefault("REGISTRY_USER", "admin")
	password = env.GetEnvOrDefault("REGISTRY_PASSWORD", "")
	mime     = env.GetEnvOrDefault("REGISTRY_MIME", MIME_V2)
)

// 💥 manifest content type should be in Conf
//
// The client should include an Accept header indicating which manifest content types it supports. For more details on the manifest format and content types, see Image Manifest Version 2, Schema 2. In a successful response, the Content-Type header will indicate which manifest type is being returned.

func NewRegistry() Registry {
	return Registry{
		// Client:     &http.Client{},
		BaseUrl: scheme + "://" + host + "/v2/",
		Conf: Conf{
			Host:     host,
			Scheme:   scheme,
			Username: username,
			Password: password,
			Mime:     mime,
		},
		HttpHeaders: map[string]string{
			"Accept":        mime,
			"Authorization": httpUtils.GetBasicAuthHeader(username, password),
		},
	}
}
