package client

import (
	"context"
	"time"
)


type Metrics struct{
	ExpiryDays time.Time
    Nameservers []string
}

// Client is a DNS client impl.
type Client interface {
	ExpireTime(ctx context.Context, domain string, host string) (Metrics, error)
}
