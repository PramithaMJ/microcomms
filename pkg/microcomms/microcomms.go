// pkg/microcomms/microcomms.go
package microcomms

import (
    "context"
    "log"
    "net/http"
    "time"
    
    "github.com/pramithamj/microcomms/internal/config"
    "github.com/pramithamj/microcomms/internal/discovery"
    "github.com/pramithamj/microcomms/internal/grpcclient"
    "github.com/pramithamj/microcomms/internal/httpclient"
    "github.com/pramithamj/microcomms/internal/mqclient"
    "github.com/rs/zerolog"
)

// Microcomms struct is the unified interface for all communication methods
type Microcomms struct {
    HTTPClient     *HTTPClient
    GRPCClient     *GRPCClient
    MQClient       *MQClient
    Discovery      *discovery.ConsulDiscovery
    CircuitBreakers map[string]*CircuitBreaker
    Logger         zerolog.Logger
    config         *config.Config
}

// HTTPClient wraps the internal HTTP client
type HTTPClient struct {
    client *httpclient.Client
}

// GRPCClient wraps the internal gRPC client
type GRPCClient struct {
    client *grpcclient.Client
}

// MQClient wraps the internal message queue client
type MQClient struct {
    queue *mqclient.MessageQueue
}

// MicrocommsConfig holds configuration for Microcomms
type MicrocommsConfig struct {
    HTTPTimeout       time.Duration
    HTTPRetryAttempts int
    ServiceDiscovery  bool
    ConsulAddress     string
    TracingEnabled    bool
    ServiceName       string
}

// DefaultConfig returns a default MicrocommsConfig
func DefaultConfig() MicrocommsConfig {
    return MicrocommsConfig{
        HTTPTimeout:       5 * time.Second,
        HTTPRetryAttempts: 3,
        ServiceDiscovery:  true,
        ConsulAddress:     "localhost:8500",
        TracingEnabled:    true,
        ServiceName:       "microcomms-client",
    }
}

// NewMicrocomms creates a new Microcomms instance with default config
func NewMicrocomms() *Microcomms {
    return NewMicrocommsWithConfig(DefaultConfig())
}

// NewMicrocommsWithConfig creates a new Microcomms instance with custom config
func NewMicrocommsWithConfig(cfg MicrocommsConfig) *Microcomms {
    // Initialize logger
    logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
    
    // Initialize tracing if enabled
    if cfg.TracingEnabled {
        InitTracing(TracingConfig{
            ServiceName:    cfg.ServiceName,
            TracingEnabled: cfg.TracingEnabled,
        })
    }
    
    // Load config from file/env
    internalCfg := config.LoadConfig()
    
    // Initialize HTTP client
    httpClient := httpclient.NewClient(httpclient.Config{
        Timeout:       cfg.HTTPTimeout,
        RetryAttempts: cfg.HTTPRetryAttempts,
    })
    
    // Initialize gRPC client
    grpcClient, err := grpcclient.NewClient("localhost:50051")
    if err != nil {
        logger.Error().Err(err).Msg("Failed to initialize gRPC client")
    }
    
    // Initialize MQ client
    mqClient := mqclient.NewMessageQueue(10)
    
    // Initialize service discovery
    var discoveryClient *discovery.ConsulDiscovery
    if cfg.ServiceDiscovery {
        dc, err := discovery.NewConsulDiscovery(cfg.ConsulAddress)
        if err != nil {
            logger.Error().Err(err).Msg("Failed to initialize service discovery")
        } else {
            discoveryClient = dc
        }
    }
    
    // Initialize circuit breakers
    circuitBreakers := make(map[string]*CircuitBreaker)
    circuitBreakers["http"] = NewCircuitBreaker("http", 5, 30*time.Second)
    circuitBreakers["grpc"] = NewCircuitBreaker("grpc", 3, 20*time.Second)
    circuitBreakers["mq"] = NewCircuitBreaker("mq", 10, 60*time.Second)
    
    return &Microcomms{
        HTTPClient: &HTTPClient{client: httpClient},
        GRPCClient: &GRPCClient{client: grpcClient},
        MQClient:   &MQClient{queue: mqClient},
        Discovery:  discoveryClient,
        CircuitBreakers: circuitBreakers,
        Logger:     logger,
        config:     internalCfg,
    }
}

// Get makes an HTTP GET request with circuit breaker and tracing
func (h *HTTPClient) Get(url string) (*http.Response, error) {
    return h.client.Get(url)
}

// GetWithContext makes an HTTP GET request with context, circuit breaker, and tracing
func (h *HTTPClient) GetWithContext(ctx context.Context, url string) (*http.Response, error) {
    ctx, span := StartSpan(ctx, "HTTPClient.GetWithContext")
    defer span.End()
    
    return h.client.Get(url)
}

// CallExample makes a gRPC call with circuit breaker and tracing
func (g *GRPCClient) CallExample(ctx context.Context, serviceMethod string) (string, error) {
    ctx, span := StartSpan(ctx, "GRPCClient.CallExample")
    defer span.End()
    
    return g.client.CallExample(ctx, serviceMethod)
}

// SendMessage sends a message to the queue with circuit breaker and tracing
func (m *MQClient) SendMessage(message string) error {
    return m.queue.SendMessage(message)
}

// SendMessageWithContext sends a message to the queue with context, circuit breaker, and tracing
func (m *MQClient) SendMessageWithContext(ctx context.Context, message string) error {
    ctx, span := StartSpan(ctx, "MQClient.SendMessageWithContext")
    defer span.End()
    
    return m.queue.SendMessage(message)
}

// ReceiveMessage receives a message from the queue
func (m *MQClient) ReceiveMessage() (string, error) {
    return m.queue.ReceiveMessage()
}

// ResolveService resolves a service using service discovery
func (m *Microcomms) ResolveService(name string) (string, error) {
    if m.Discovery == nil {
        return "", ErrServiceDiscoveryNotEnabled
    }
    
    return m.Discovery.FindService(name)
}

// Get makes an HTTP GET request to a service using service discovery
func (m *Microcomms) Get(ctx context.Context, serviceName, path string) (*http.Response, error) {
    ctx, span := StartSpan(ctx, "Microcomms.Get")
    defer span.End()
    
    // Circuit breaker pattern
    var resp *http.Response
    var err error
    
    err = m.CircuitBreakers["http"].Execute(func() error {
        // Resolve service name if using service discovery
        serviceURL := serviceName
        if m.Discovery != nil {
            resolvedURL, resolveErr := m.Discovery.FindService(serviceName)
            if resolveErr != nil {
                return resolveErr
            }
            serviceURL = resolvedURL
        }
        
        // Make the request
        resp, err = m.HTTPClient.GetWithContext(ctx, serviceURL+path)
        return err
    })
    
    return resp, err
}