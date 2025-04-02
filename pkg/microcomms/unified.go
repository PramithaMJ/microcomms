package microcomms

import (
    "context"
    "fmt"
    "time"
)

// MessageRequest represents a generic communication request
type MessageRequest struct {
    Target  string            // Service name or URL
    Payload interface{}       // Message payload
    Headers map[string]string // Headers for HTTP/gRPC
    Timeout time.Duration     // Request timeout
}

// MessageResponse represents a generic communication response
type MessageResponse struct {
    StatusCode int         // HTTP status code or gRPC status
    Payload    interface{} // Response payload
    Headers    map[string]string
    Protocol   string // Which protocol was used (HTTP, gRPC, MQ)
}

// ProtocolType defines the communication protocol to use
type ProtocolType string

const (
    ProtocolHTTP    ProtocolType = "http"
    ProtocolGRPC    ProtocolType = "grpc"
    ProtocolMQ      ProtocolType = "mq"
    ProtocolAuto    ProtocolType = "auto" // Auto-select best protocol
    ProtocolFallback ProtocolType = "fallback" // Try protocols in order (http -> grpc -> mq)
)

// Send sends a message using the specified protocol
// If protocol is Auto, it will select the best protocol based on the message type
// If protocol is Fallback, it will try HTTP first, then gRPC, then MQ
func (m *Microcomms) Send(ctx context.Context, req MessageRequest, protocol ProtocolType) (*MessageResponse, error) {
    switch protocol {
    case ProtocolHTTP:
        return m.sendHTTP(ctx, req)
    case ProtocolGRPC:
        return m.sendGRPC(ctx, req)
    case ProtocolMQ:
        return m.sendMQ(ctx, req)
    case ProtocolFallback:
        return m.sendWithFallback(ctx, req)
    case ProtocolAuto:
        return m.sendAuto(ctx, req)
    default:
        return nil, fmt.Errorf("unsupported protocol: %s", protocol)
    }
}

// sendHTTP sends a message over HTTP
func (m *Microcomms) sendHTTP(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
    // Implementation details here
    // For now, just a simple wrapper over the existing HTTP client
    resp, err := m.HTTPClient.Get(req.Target)
    if err != nil {
        return nil, err
    }
    
    return &MessageResponse{
        StatusCode: resp.StatusCode,
        Protocol:   "http",
        Headers:    make(map[string]string),
    }, nil
}

// sendGRPC sends a message over gRPC
func (m *Microcomms) sendGRPC(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
    // Implementation details here
    resp, err := m.GRPCClient.CallExample(ctx, req.Target)
    if err != nil {
        return nil, err
    }
    
    return &MessageResponse{
        Payload:  resp,
        Protocol: "grpc",
    }, nil
}

// sendMQ sends a message through a message queue
func (m *Microcomms) sendMQ(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
    // Implementation details here
    payload, ok := req.Payload.(string)
    if !ok {
        return nil, fmt.Errorf("MQ payload must be a string")
    }
    
    err := m.MQClient.SendMessage(payload)
    if err != nil {
        return nil, err
    }
    
    return &MessageResponse{
        Protocol: "mq",
    }, nil
}

// sendWithFallback tries HTTP first, then gRPC, then MQ
func (m *Microcomms) sendWithFallback(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
    resp, err := m.sendHTTP(ctx, req)
    if err == nil {
        return resp, nil
    }
    
    resp, err = m.sendGRPC(ctx, req)
    if err == nil {
        return resp, nil
    }
    
    return m.sendMQ(ctx, req)
}

// sendAuto selects the best protocol based on the message type
func (m *Microcomms) sendAuto(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
    // Simple heuristic: if target contains http/https, use HTTP
    if len(req.Target) > 4 && (req.Target[:4] == "http" || req.Target[:5] == "https") {
        return m.sendHTTP(ctx, req)
    }
    
    // If target contains a dot (likely a service name), use gRPC
    for _, c := range req.Target {
        if c == '.' {
            return m.sendGRPC(ctx, req)
        }
    }
    
    // Default to MQ for simple string messages
    return m.sendMQ(ctx, req)
}