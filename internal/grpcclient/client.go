package grpcclient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

// Client struct
type Client struct {
	conn *grpc.ClientConn
}

// NewClient initializes a new gRPC client
func NewClient(target string) (*Client, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	return &Client{conn: conn}, nil
}

// CallExample makes a sample gRPC call
func (c *Client) CallExample(ctx context.Context, serviceMethod string) (string, error) {
	// This is where you can add your gRPC method call logic
	time.Sleep(time.Second) // Simulate delay
	return "Response from " + serviceMethod, nil
}
