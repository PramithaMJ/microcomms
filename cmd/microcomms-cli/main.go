package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/pramithamj/microcomms/pkg/microcomms"
)

func main() {
    // Initialize the microcomms client
    client := microcomms.NewMicrocomms()
    
    // Example of HTTP request
    fmt.Println("Making HTTP request...")
    resp, err := client.HTTPClient.Get("https://jsonplaceholder.typicode.com/posts/1")
    if err != nil {
        log.Fatalf("HTTP request failed: %v", err)
    }
    fmt.Println("HTTP Response Status:", resp.Status)
    
    // Example of gRPC call
    fmt.Println("\nMaking gRPC call...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    result, err := client.GRPCClient.CallExample(ctx, "ExampleService.Method")
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
    } else {
        fmt.Println("gRPC Response:", result)
    }
    
    // Example of message queue
    fmt.Println("\nSending message to queue...")
    err = client.MQClient.SendMessage("Hello from Microcomms!")
    if err != nil {
        log.Printf("Failed to send message: %v", err)
    }
    
    // Receive message from queue
    fmt.Println("\nReceiving message from queue...")
    message, err := client.MQClient.ReceiveMessage()
    if err != nil {
        log.Printf("Failed to receive message: %v", err)
    } else {
        fmt.Println("Received message:", message)
    }
    
    fmt.Println("\nMicrocomms demonstration completed!")
}