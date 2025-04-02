package microcomms

import (
    "errors"
    "fmt"
)

// Common errors
var (
    ErrServiceUnavailable       = errors.New("service unavailable")
    ErrCircuitBreakerOpen       = errors.New("circuit breaker open")
    ErrRateLimited              = errors.New("rate limited")
    ErrServiceDiscoveryNotEnabled = errors.New("service discovery not enabled")
    ErrTimeout                  = errors.New("request timed out")
    ErrInvalidProtocol          = errors.New("invalid protocol")
)

// ServiceError represents an error from a service
type ServiceError struct {
    ServiceName string
    StatusCode  int
    Message     string
    Err         error
}

func (e *ServiceError) Error() string {
    return fmt.Sprintf("service %s error (status %d): %s", e.ServiceName, e.StatusCode, e.Message)
}

func (e *ServiceError) Unwrap() error {
    return e.Err
}

// NewServiceError creates a new ServiceError
func NewServiceError(serviceName string, statusCode int, message string, err error) *ServiceError {
    return &ServiceError{
        ServiceName: serviceName,
        StatusCode:  statusCode,
        Message:     message,
        Err:         err,
    }
}