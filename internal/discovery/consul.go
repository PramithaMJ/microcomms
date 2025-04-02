package discovery

import (
    "fmt"
    "sync"
    "time"
    
    "github.com/hashicorp/consul/api"
)

// ConsulDiscovery implements service discovery using Consul
type ConsulDiscovery struct {
    client        *api.Client
    serviceCache  map[string][]*api.ServiceEntry
    cacheMutex    sync.RWMutex
    cacheExpiry   time.Duration
    lastCacheTime map[string]time.Time
}

// NewConsulDiscovery creates a new Consul discovery client
func NewConsulDiscovery(addr string) (*ConsulDiscovery, error) {
    config := api.DefaultConfig()
    if addr != "" {
        config.Address = addr
    }
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create Consul client: %v", err)
    }
    
    return &ConsulDiscovery{
        client:        client,
        serviceCache:  make(map[string][]*api.ServiceEntry),
        cacheExpiry:   time.Second * 30, // Cache for 30 seconds
        lastCacheTime: make(map[string]time.Time),
    }, nil
}

// FindService finds a service by name
func (c *ConsulDiscovery) FindService(name string) (string, error) {
    entries, err := c.getServiceEntries(name)
    if err != nil {
        return "", err
    }
    
    if len(entries) == 0 {
        return "", fmt.Errorf("no instances of service '%s' found", name)
    }
    
    // For simplicity, return the first instance
    // In a real implementation, we could implement load balancing here
    instance := entries[0]
    return fmt.Sprintf("http://%s:%d", instance.Service.Address, instance.Service.Port), nil
}

// getServiceEntries gets service entries, using cache if available
func (c *ConsulDiscovery) getServiceEntries(name string) ([]*api.ServiceEntry, error) {
    c.cacheMutex.RLock()
    entries, exists := c.serviceCache[name]
    lastCacheTime := c.lastCacheTime[name]
    c.cacheMutex.RUnlock()
    
    // Return from cache if it's still valid
    if exists && time.Since(lastCacheTime) < c.cacheExpiry {
        return entries, nil
    }
    
    // Otherwise query Consul
    c.cacheMutex.Lock()
    defer c.cacheMutex.Unlock()
    
    // Check again after acquiring the lock
    if entries, exists := c.serviceCache[name]; exists && time.Since(c.lastCacheTime[name]) < c.cacheExpiry {
        return entries, nil
    }
    
    entries, _, err := c.client.Health().Service(name, "", true, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to query Consul for service '%s': %v", name, err)
    }
    
    c.serviceCache[name] = entries
    c.lastCacheTime[name] = time.Now()
    
    return entries, nil
}

// RegisterService registers a service with Consul
func (c *ConsulDiscovery) RegisterService(name, address string, port int) error {
    service := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s-%d", name, address, port),
        Name:    name,
        Address: address,
        Port:    port,
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
            Interval: "10s",
            Timeout:  "1s",
        },
    }
    
    return c.client.Agent().ServiceRegister(service)
}

// DeregisterService deregisters a service from Consul
func (c *ConsulDiscovery) DeregisterService(serviceID string) error {
    return c.client.Agent().ServiceDeregister(serviceID)
}