package resolve

import (
	"context"
	"net"
	"time"
)

// NewResolver creates a new DNS resolver that uses a custom DNS server address.
// The resolver will prefer the Go DNS resolver over the system resolver.
//
// Parameters:
// - resolverAddr: The address of the DNS server to use for resolution.
//
// Returns:
// - A pointer to a net.Resolver configured to use the specified DNS server.
func NewResolver(resolverAddr string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, resolverAddr)
		},
	}
}
