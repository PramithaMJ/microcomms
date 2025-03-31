package microcomms

import (
    "log"
    "net/http"
    "time"
    "context"
    
    "github.com/pramithamj/microcomms/internal/httpclient"
    "github.com/pramithamj/microcomms/internal/grpcclient"
    "github.com/pramithamj/microcomms/internal/mqclient"
    "github.com/pramithamj/microcomms/internal/config"
)

// Microcomms struct is the unified interface for all communication methods
type Microcomms struct {
    HTTPClient *HTTPClient
    GRPCClient *GRPCClient
    MQClient   *MQClient
    config     *config.Config
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

// NewMicrocomms creates a new Microcomms instance
func NewMicrocomms() *Microcomms {
    cfg := config.LoadConfig()

    httpClient := httpclient.NewClient(httpclient.Config{
        Timeout:       cfg.Timeout,
        RetryAttempts: cfg.RetryAttempts,
    })

    grpcClient, err := grpcclient.NewClient("localhost:50051")
    if err != nil {
        log.Fatalf("Error initializing gRPC client: %v", err)
    }

    mqClient := mqclient.NewMessageQueue(10)

    return &Microcomms{
        HTTPClient: &HTTPClient{client: httpClient},
        GRPCClient: &GRPCClient{client: grpcClient},
        MQClient:   &MQClient{queue: mqClient},
        config:     cfg,
    }
}

// Get makes an HTTP GET request
func (h *HTTPClient) Get(url string) (*http.Response, error) {
    return h.client.Get(url)
}

// CallExample makes a gRPC call
func (g *GRPCClient) CallExample(ctx context.Context, serviceMethod string) (string, error) {
    return g.client.CallExample(ctx, serviceMethod)
}

// SendMessage sends a message to the queue
func (m *MQClient) SendMessage(message string) error {
    return m.queue.SendMessage(message)
}

// ReceiveMessage receives a message from the queue
func (m *MQClient) ReceiveMessage() (string, error) {
    return m.queue.ReceiveMessage()
}