package dnsresolver

import (
	"context"
	"errors"
	"net"
)

type Resolver struct {
	resolver *net.Resolver
}

func NewResolver() *Resolver {
	return &Resolver{resolver: net.DefaultResolver}
}

func (r *Resolver) LookupIP(ctx context.Context, host string) ([]net.IP, error) {
	if host == "" {
		return nil, errors.New("host is empty")
	}

	return r.resolver.LookupIP(ctx, "ip", host)
}
