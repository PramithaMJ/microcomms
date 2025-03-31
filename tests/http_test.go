package tests

import (
	"testing"
	"github.com/pramithamj/microcomms/internal/httpclient"
)

func TestHTTPClient_Get(t *testing.T) {
	client := httpclient.NewClient(httpclient.Config{
		Timeout:       5,
		RetryAttempts: 2,
	})

	resp, err := client.Get("http://example.com")
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if resp == nil {
		t.Fatalf("Expected response, but got nil")
	}
}
