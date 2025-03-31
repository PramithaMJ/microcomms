package discovery

import (
	"fmt"
	"log"
)

// DiscoverService finds the service in the registry
func DiscoverService(serviceName string) (string, error) {
	// This is a placeholder. You can use Consul, Etcd, or Kubernetes DNS.
	// Assuming the service is found at http://localhost:8080
	return fmt.Sprintf("http://localhost:8080/%s", serviceName), nil
}

// ServiceResolver resolves the service URL dynamically
func ServiceResolver(serviceName string) string {
	url, err := DiscoverService(serviceName)
	if err != nil {
		log.Fatalf("Failed to resolve service: %v", err)
	}
	return url
}
