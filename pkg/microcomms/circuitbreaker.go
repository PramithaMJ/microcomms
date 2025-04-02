package microcomms

import (
    "fmt"
    "sync"
    "time"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
    StateClosed CircuitBreakerState = iota // Normal operation, requests pass through
    StateOpen                              // Circuit is open, requests fail fast
    StateHalfOpen                          // Testing if circuit can be closed again
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
    name          string
    state         CircuitBreakerState
    failureCount  int
    failureThreshold int
    resetTimeout  time.Duration
    lastFailureTime time.Time
    mutex         sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        name:            name,
        state:           StateClosed,
        failureThreshold: failureThreshold,
        resetTimeout:    resetTimeout,
    }
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.AllowRequest() {
        return fmt.Errorf("circuit breaker '%s' is open", cb.name)
    }
    
    err := fn()
    
    if err != nil {
        cb.RecordFailure()
        return err
    }
    
    cb.RecordSuccess()
    return nil
}

// AllowRequest checks if a request is allowed to pass through the circuit breaker
func (cb *CircuitBreaker) AllowRequest() bool {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        // Check if reset timeout has expired
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            // Allow one request to test the circuit
            cb.mutex.RUnlock()
            cb.mutex.Lock()
            cb.state = StateHalfOpen
            cb.mutex.Unlock()
            cb.mutex.RLock()
            return true
        }
        return false
    case StateHalfOpen:
        // In half-open state, allow only one request
        return true
    default:
        return false
    }
}

// RecordFailure records a failure and potentially opens the circuit
func (cb *CircuitBreaker) RecordFailure() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    cb.failureCount++
    cb.lastFailureTime = time.Now()
    
    if cb.state == StateHalfOpen || (cb.state == StateClosed && cb.failureCount >= cb.failureThreshold) {
        cb.state = StateOpen
    }
}

// RecordSuccess records a success and potentially closes the circuit
func (cb *CircuitBreaker) RecordSuccess() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    if cb.state == StateHalfOpen {
        cb.failureCount = 0
        cb.state = StateClosed
    } else if cb.state == StateClosed {
        cb.failureCount = 0
    }
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitBreakerState {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    return cb.state
}