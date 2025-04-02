package microcomms

import (
    "sync"
    "time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
    tokensPerSecond float64
    maxTokens       float64
    tokens          float64
    lastRefill      time.Time
    mutex           sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(tokensPerSecond, maxTokens float64) *RateLimiter {
    return &RateLimiter{
        tokensPerSecond: tokensPerSecond,
        maxTokens:       maxTokens,
        tokens:          maxTokens,
        lastRefill:      time.Now(),
    }
}

// Allow checks if a request is allowed by the rate limiter
func (rl *RateLimiter) Allow() bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    // Refill tokens
    now := time.Now()
    elapsed := now.Sub(rl.lastRefill).Seconds()
    rl.lastRefill = now
    
    // Add tokens based on elapsed time
    rl.tokens += elapsed * rl.tokensPerSecond
    if rl.tokens > rl.maxTokens {
        rl.tokens = rl.maxTokens
    }
    
    // Check if we have enough tokens
    if rl.tokens < 1 {
        return false
    }
    
    // Use a token
    rl.tokens--
    return true
}

// WaitAndAllow waits until a token is available and then uses it
func (rl *RateLimiter) WaitAndAllow() {
    for {
        if rl.Allow() {
            return
        }
        time.Sleep(time.Millisecond * 10)
    }
}

// RateLimitedFunc wraps a function with rate limiting
func (m *Microcomms) RateLimitedFunc(limiter *RateLimiter, fn func() error) error {
    if !limiter.Allow() {
        return ErrRateLimited
    }
    return fn()
}