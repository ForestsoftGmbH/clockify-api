package clockify_api

import (
	"net/http"
	"os"
)

type Api struct {
	HttpClient  HTTPClient
	Credentials Credentials
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Credentials struct {
	WorkspaceId string
	ApiKey      string
}

func NewApi() *Api {
	return &Api{
		HttpClient:  http.DefaultClient,
		Credentials: Credentials{WorkspaceId: os.Getenv("CLOCKIFY_WORKSPACE"), ApiKey: os.Getenv("CLOCKIFY_API_KEY")},
	}
}
