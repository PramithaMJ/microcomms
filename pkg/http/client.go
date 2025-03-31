package http

import "github.com/pramithamj/microcomms/internal/httpclient"

func NewClient(baseURL string) *httpclient.Client {
	return httpclient.NewClient(baseURL)
}
