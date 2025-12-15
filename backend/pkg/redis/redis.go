package redis

import "context"

// NewClient is a stub for smoke tests. In production replace with a
// real Redis client (github.com/redis/go-redis/v11).
func NewClient(addr string) interface{} {
    _ = addr
    _ = context.Background()
    return nil
}
